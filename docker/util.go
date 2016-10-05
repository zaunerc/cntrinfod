package docker

import (
	"fmt"
	"github.com/docker/engine-api/client"
	"sync"
)

var dockerClients = make(map[string]*client.Client)
var mutex = &sync.Mutex{}

// getDockerClientForUrl returns a Docker HTTP client specific
// for one URL. The HTTP client will pool and reuse idle
// connections to Docker.
func GetDockerClientForUrl(dockerUrl string, version string) (*client.Client, error) {

	mutex.Lock()
	defer mutex.Unlock()

	if dockerClients[dockerUrl] == nil {

		var err error

		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		dockerClients[dockerUrl], err = client.NewClient(dockerUrl, version, nil, defaultHeaders)

		if err != nil {
			// E.g. only root may connect to docker socket.
			fmt.Printf("Error while creating docker HTTP client: %s", err)
			return nil, err
		}
	}

	return dockerClients[dockerUrl], nil
}
