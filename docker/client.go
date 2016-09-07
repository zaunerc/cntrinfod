/*
- https://godoc.org/github.com/docker/engine-api/client#Client.ContainerInspect
*/

package docker

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/docker/engine-api/client"
	"golang.org/x/net/context"
)

func FetchHostHostname() string {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		// Only root may connect to docker socket.
		panic(err)
	}

	info, _ := cli.Info(context.Background())

	return info.Name
}

func FetchHostInfo() string {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		// Only root may connect to docker socket.
		panic(err)
	}

	info, _ := cli.Info(context.Background())
	infoString := spew.Sdump(info)

	return infoString
}
