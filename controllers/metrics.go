package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	ormMetric "gorm.io/plugin/prometheus"
)

var (
	GaugeDomainCounts = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "recorednsui",
		Subsystem: "domains",
		Name:      "count",
		Help:      "domains managed in reCoreDNS-UI",
	})

	GaugeRecordCounts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "recorednsui",
		Subsystem: "records",
		Name:      "count",
		Help:      "records managed in reCoreDNS-UI, by domain",
	}, []string{"domain"})
)

func RegisterMetrics() {
	prometheus.MustRegister(GaugeDomainCounts, GaugeRecordCounts)

	GormMetrics := ormMetric.New(ormMetric.Config{
		DBName: "recoredns-ui",
		MetricsCollector: []ormMetric.MetricsCollector{
			&ormMetric.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		},
	}).Collectors
	prometheus.MustRegister(GormMetrics...)

	GinMetrics := ginprometheus.NewPrometheus("recorednsui")
	for _, v := range GinMetrics.MetricsList {
		prometheus.Register(v.MetricCollector)
	}
}

func RefreshMetrics() error {
	domainCounts, err := getDomainCounts()
	if err != nil {
		return err
	}
	GaugeDomainCounts.Set(domainCounts)

	recordCounts, err := getRecordCounts()
	if err != nil {
		return err
	}

	for domain, counts := range recordCounts {
		GaugeRecordCounts.WithLabelValues(domain).Set(counts)
	}

	return nil
}
