package api

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func RegisterNode(port, peersArg string) {
	nodeID := "localhost:" + port
	n := &Node{
		ID:     nodeID,
		Addr:   "http://" + nodeID,
		Peers:  make(map[string]*Peer),
		Client: &http.Client{Timeout: 2 * time.Second},
	}

	http.HandleFunc("/join", n.joinHandler)
	http.HandleFunc("/peers", n.peersHandler)
	http.HandleFunc("/heartbeat", n.heartbeatHandler)

	// Counter APIs - later we can move to other pkg
	http.HandleFunc("/counter/increment", n.IncrementHandler)
	http.HandleFunc("/counter/decrement", n.DecrementHandler)
	http.HandleFunc("/counter/value", n.ValueHandler)

	// sync from other node
	http.HandleFunc("/counter/sync", n.SyncHandler)

	if peersArg != "" {
		for _, addr := range strings.Split(peersArg, ",") {
			n.bootstrap(addr)
		}
	}

	// Start heartbeat loop
	go n.heartbeatLoop()

	// Start server
	log.Printf("[%s] Starting node on port %s", n.ID, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
