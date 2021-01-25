package helpers

import (
	"log"

	"github.com/phayes/freeport"
)

// GetOpenedPort : return opened TCP port
func GetOpenedPort() int {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	return port
}
