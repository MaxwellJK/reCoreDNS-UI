package controllers

import (
	"fmt"
	"reCoreD-UI/database"
	"reCoreD-UI/models"
	"strconv"

	"github.com/sirupsen/logrus"

	"reflect"
)

type domainsDAO struct {
	database.BaseDAO[models.IDomain]
}

func CreateDomain(d *models.Domain) (*models.Domain, error) {
	nss, err := GetDNS()
	if err != nil {
		return nil, err
	}

	tx := database.Client.Begin()
	if _, err := (domainsDAO{}).Create(tx, d); err != nil {
		tx.Rollback()
		return nil, err
	}

	r := &models.Record[models.SOARecord]{}
	r.Zone = d.WithDotEnd()
	r.Name = "@"
	r.RecordType = models.RecordTypeSOA
	r.Content = d.GenerateSOA()
	logrus.Debug(r)
	if err := r.CheckZone(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := (recordsDAO{}).Create(tx, r); err != nil {
		tx.Rollback()
		return nil, err
	}

	for i, ns := range nss {
		record := &models.Record[models.NSRecord]{
			Zone:       d.WithDotEnd(),
			RecordType: models.RecordTypeNS,
			Name:       fmt.Sprintf("ns%d", i+1),
		}
		record.Content.Host = ns

		if _, err := (recordsDAO{}).Create(tx, record); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return d, tx.Commit().Error
}

func GetDomains(domain string) ([]models.Domain, error) {
	if domain != "" {
		fmt.Print("getDomains when a domain is provided")
		r, err := (domainsDAO{}).GetAll(database.Client, &models.Domain{DomainName: domain}, &models.Domain{DomainName: domain})
		n := make([]models.Domain, 0)
		for _, e := range r {
			i, ok := e.(*models.Domain)
			if !ok {
				continue
			}
			n = append(n, *i)
		}
		return n, err
	} else {
		r, err := (domainsDAO{}).GetAll(database.Client, &models.Domain{}, &models.Domain{})
		n := make([]models.Domain, 0)
		for _, e := range r {
			i, ok := e.(*models.Domain)
			if !ok {
				continue
			}
			n = append(n, *i)
		}
		return n, err
	}
}

func GetDomainByID(domainId uint) (models.Domain, error) {
	r, err := (domainsDAO{}).GetAll(database.Client, &models.Domain{ID: domainId})

	n := make([]models.Domain, 0)
	for _, e := range r {
		i, ok := e.(*models.Domain)
		if !ok {
			continue
		}
		n = append(n, *i)
	}
	return n[0], err
}

func structToMap(s interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(s)

	// Check if the input is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	// Iterate over the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)
		fmt.Println(field)
		fmt.Println(value)

		// Get the JSON tag (if any)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			// If no JSON tag exists, use the field name
			jsonTag = field.Name
		}

		// Store the field's value in the map using the JSON tag as the key
		result[jsonTag] = value.Interface()
	}

	return result, nil
}

func UpdateDomain(d *models.Domain) error {
	tx := database.Client.Begin()

	oldDomain, err := GetDomainByID(d.ID)

	if err != nil {
		return err
	}

	oldDomainName := oldDomain.WithDotEnd()

	records, err := GetRecordsMF(oldDomainName)
	if err != nil {
		return err
	}

	var iRecords []models.IRecord
	for _, v := range records {
		if v.RecordType != "SOA" {
			v.Zone = d.WithDotEnd()
		} else {

			soa := d.GenerateSOA()
			soaMap, err := structToMap(soa.SOARecord)
			if err != nil {
				tx.Rollback()
				return err
			}
			fmt.Println(soaMap)
			v.Content = soaMap
			v.Zone = d.WithDotEnd()
			if err := v.CheckZone(); err != nil {
				tx.Rollback()
				return err
			}

		}

		iRecords = append(iRecords, &v)
	}

	for _, r := range iRecords {
		err = UpdateRecord(r)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	logrus.Debug(d)
	if _, err := (domainsDAO{}).Update(tx, d); err != nil {
		tx.Rollback()
		return err
	}

	// soa, err := (recordsDAO{}).GetOne(tx, &models.Record[models.RecordContentDefault]{RecordType: models.RecordTypeSOA}, &models.Record[models.RecordContentDefault]{Zone: oldDomainName})
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// fmt.Print("this is the soa ")
	// fmt.Print(soa)

	// r := &models.Record[models.SOARecord]{}
	// if err := r.FromEntity(soa); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// r.Content = d.GenerateSOA()
	// r.Zone = d.WithDotEnd()
	// if err := r.CheckZone(); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// if _, err := (recordsDAO{}).Update(tx, r); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	fmt.Print("End of the uxpdate domain...")
	return tx.Commit().Error
}

func DeleteDomain(id string) error {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	tx := database.Client.Begin()
	domain, err := (domainsDAO{}).GetOne(tx, &models.Domain{ID: uint(ID)})
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := (domainsDAO{}).Delete(tx, &models.Domain{ID: uint(ID)}); err != nil {
		tx.Rollback()
		return err
	}

	if err := (recordsDAO{}).Delete(tx, &models.Record[models.RecordContentDefault]{}, &models.Record[models.RecordContentDefault]{Zone: domain.WithDotEnd()}); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// for metrics
func getDomainCounts() (float64, error) {
	c, err := (domainsDAO{}).GetAll(database.Client, &models.Domain{})
	if err != nil {
		return 0, err
	}
	return float64(len(c)), nil
}
