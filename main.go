package main

import (
	"log"
	"os"
	"strings"

	"github.com/Shubhacker/asset-discovery/api"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go --port=8081 [--peers=host1:port1,host2:port2]")
	}

	var port string
	var peersArg string
	for _, arg := range os.Args[1:] {
		// get port/peers
		if strings.HasPrefix(arg, "--port=") {
			port = strings.TrimPrefix(arg, "--port=")
		} else if strings.HasPrefix(arg, "--peers=") {
			peersArg = strings.TrimPrefix(arg, "--peers=")
		}
	}
	if port == "" {
		log.Fatal("must provide --port")
	}
	api.RegisterNode(port, peersArg)
}
