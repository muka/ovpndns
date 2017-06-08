package main

import (
	"os"

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
	}

	app.Action = func(c *cli.Context) error {

		src := c.String("src")
		out := c.String("out")
		domain := c.String("domain")

		WatchFile(file)

		return nil
	}

	app.Run(os.Args)

}
