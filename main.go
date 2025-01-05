package main

import (
	"os"
	"reCoreDNS-UI/cmd/config"
	"reCoreDNS-UI/cmd/server"
	_ "reCoreDNS-UI/docs"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func init() {
	logrus.SetReportCaller(true)
}

// @title						reCoreDNS-UI API
// @version					1.0
// @description				APIs for reCoreDNS-UI
// @BasePath					/api/v1
// @securityDefinitions.basic	BasicAuth
func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "config yaml file",
			Aliases: []string{"c"},
			EnvVars: []string{"RECOREDNS_CONFIG_FILE"},
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "mysql-dsn",
			Usage:   "mysql dsn",
			EnvVars: []string{"RECOREDNS_MYSQL_DSN"},
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug mode",
			Value: false,
			Action: func(ctx *cli.Context, b bool) error {
				if b {
					logrus.SetLevel(logrus.DebugLevel)
				}
				return nil
			},
		}),
	}

	app := &cli.App{
		Name:  "reCoreDNS-UI",
		Usage: "Web UI for CoreDNS",
		Before: altsrc.InitInputSourceWithContext(
			flags, altsrc.NewYamlSourceFromFlagFunc("config"),
		),
		Flags: flags,
		Commands: []*cli.Command{
			server.Command,
			config.Command,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
