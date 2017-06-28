package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/muka/ovpndns/ddns"
	"github.com/muka/ovpndns/parser"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()

	app.Name = "ovpndns"
	app.Usage = "Parse a openvpn status file and store to hosts like file for dnsmasq or push to ddns"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "src, s",
			Value:  "./tmp/openvpn-status.log",
			Usage:  "Set the openvpn status file source",
			EnvVar: "OVPN_STATUS_FILE",
		},
		cli.StringFlag{
			Name:   "out, o",
			Value:  "",
			Usage:  "Set the output file of a hosts-like formatted list of clients",
			EnvVar: "OUT_FILE",
		},
		cli.StringFlag{
			Name:   "domain, d",
			Value:  "service.local",
			Usage:  "Set the default domain to append to each host name",
			EnvVar: "DOMAIN",
		},
		cli.StringFlag{
			Name:   "ddns",
			Value:  "127.0.0.1:5551",
			Usage:  "DDNS API host",
			EnvVar: "DDNS_HOST",
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
		ddnsHost := c.String("ddns")

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if ddnsHost != "" {
			ddns.CreateClient(ddnsHost)
		}

		go func() {
			updates := parser.GetChannel()
			for {
				select {
				case records := <-updates:

					log.Debugf("Updating configuration, %d records", len(records))

					var buffer bytes.Buffer
					for i, record := range records {

						records[i].Name = record.Name + "." + domain
						c := fmt.Sprintf("%s %s.%s", record.IP, record.Name, domain)

						log.Debugf("Add line %s", c)
						buffer.WriteString(c)
						buffer.WriteString("\n")
					}
					log.Debugf("Storing to file %s", out)
					ioutil.WriteFile(out, buffer.Bytes(), 0644)

					if ddnsHost != "" {
						ddns.Compare(records)
					}
				}
			}
		}()

		go parser.ParseFile(src)
		parser.WatchFile(src)

		return nil
	}

	app.Run(os.Args)

}
