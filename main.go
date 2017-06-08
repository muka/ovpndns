package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/muka/ovpndns/parser"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()

	app.Name = "ovpndns"
	app.Usage = "Parse openvpn status file and push hosts to dnsmasq"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "src, s",
			Value:  "./tmp/openvpn-status.log",
			Usage:  "Set the openvpn status file source",
			EnvVar: "OVPN_STATUS_FILE",
		},
		cli.StringFlag{
			Name:   "out, o",
			Value:  "./data/hosts",
			Usage:  "Set the dns record output file",
			EnvVar: "OUT_FILE",
		},
		cli.StringFlag{
			Name:   "domain, d",
			Value:  "service.local",
			Usage:  "Set the default domain",
			EnvVar: "DOMAIN",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debugging logs",
			EnvVar: "DEBUG",
		},
	}

	app.Action = func(c *cli.Context) error {

		debug := c.Bool("debug")
		src := c.String("src")
		out := c.String("out")
		domain := c.String("domain")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		go func() {
			updates := parser.GetChannel()
			for {
				select {
				case <-updates:

					log.Debug("Updating configuration")

					var buffer bytes.Buffer
					records := parser.GetRecords()
					for _, record := range records {
						c := fmt.Sprintf("address=/%s.%s/%s", record.Name, domain, record.IP)
						log.Debugf("Add line %s", c)
						buffer.WriteString(c)
						buffer.WriteString("\n")
					}
					log.Debugf("Storing to file %s", out)
					ioutil.WriteFile(out, buffer.Bytes(), 0644)
				}
			}
		}()

		go parser.ParseFile(src)
		parser.WatchFile(src)

		return nil
	}

	app.Run(os.Args)

}
