/*
- https://godoc.org/github.com/docker/engine-api/client#Client.ContainerInspect
*/

package docker

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
)

// FetchHostName returns the hostname of the Docker
// host OS. In case of any errors an empty string
// is returned.
func FetchHostHostname() string {

	cli, err := GetDockerClientForUrl("unix:///var/run/docker.sock", "v1.22")

	if err != nil {
		fmt.Printf("Error while trying to get docker HTTP client: %s", err)
		return ""
	}

	info, err := cli.Info(context.Background())

	if err != nil {
		fmt.Printf("Error while calling docker HTTP API: %s", err)
		return ""
	} else {
		return info.Name
	}
}

// FetchHostInfo returns the Docker servers
// https://godoc.org/github.com/docker/engine-api/client#Client.Info data
// structure. It is converted to a string and formatted in
// a nice way by go-spew.
func FetchHostInfo() string {

	cli, err := GetDockerClientForUrl("unix:///var/run/docker.sock", "v1.22")

	if err != nil {
		fmt.Printf("Error while trying to get docker HTTP client: %s", err)
		return ""
	}

	info, err := cli.Info(context.Background())

	if err != nil {
		fmt.Printf("Error while calling docker HTTP API: %s", err)
		return ""
	} else {
		infoString := spew.Sdump(info)
		return infoString
	}

}
