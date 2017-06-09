package docker

import (
	"github.com/docker/docker/client"
)

var dockerCli *client.Client

//GetClient return a docker client instance
func GetClient() (*client.Client, error) {

	if dockerCli != nil {
		return dockerCli, nil
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	dockerCli = cli

	return cli, nil
}

//SendSIGHUP to a dnsmasq container to force config reloading
func SendSIGHUP() error {

	// cli, err := GetClient()
	// if err != nil {
	// 	return err
	// }

	// ctx := context.Background()
	return nil
}
