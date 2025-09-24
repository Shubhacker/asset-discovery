package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Peer info
type Peer struct {
	ID       string    `json:"id"`
	Address  string    `json:"address"`
	LastSeen time.Time `json:"last_seen"`
	Alive    bool      `json:"alive"`
}

// Discovery state
type Node struct {
	ID      string
	Addr    string
	Peers   map[string]*Peer
	m       sync.RWMutex
	Client  *http.Client
	Counter int
}

// Create new node
func (n *Node) joinHandler(w http.ResponseWriter, r *http.Request) {
	var newPeer Peer
	if err := json.NewDecoder(r.Body).Decode(&newPeer); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	n.m.Lock()
	defer n.m.Unlock()

	if _, exists := n.Peers[newPeer.ID]; !exists {
		newPeer.LastSeen = time.Now()
		newPeer.Alive = true
		n.Peers[newPeer.ID] = &newPeer
		log.Printf("[%s] Added new peer: %s", n.ID, newPeer.ID)
	}

	var peers []*Peer
	for _, p := range n.Peers {
		peers = append(peers, p)
	}

	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"peers": peers,
	})
}

// Get peers
func (n *Node) peersHandler(w http.ResponseWriter, r *http.Request) {
	n.m.RLock()
	defer n.m.RUnlock()

	var peers []*Peer
	for _, p := range n.Peers {
		peers = append(peers, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"peers": peers,
	})
}

// health check, later we can add more validations
func (n *Node) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (n *Node) bootstrap(peerAddr string) {
	url := fmt.Sprintf("http://%s/join", peerAddr)
	reqBody := map[string]string{
		"id":      n.ID,
		"address": n.Addr,
	}
	body, _ := json.Marshal(reqBody)

	resp, err := n.Client.Post(url, "application/json", strings.NewReader(string(body)))
	if err != nil {
		log.Printf("[%s] Failed to join %s: %v", n.ID, peerAddr, err)
		return
	}
	defer resp.Body.Close()

	var data map[string][]*Peer
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("[%s] Failed to decode join response: %v", n.ID, err)
		return
	}

	// Merge peers
	n.m.Lock()
	defer n.m.Unlock()
	for _, p := range data["peers"] {
		if p.ID != n.ID {
			n.Peers[p.ID] = p
		}
	}
	log.Printf("[%s] Bootstrapped with %d peers", n.ID, len(n.Peers))
}

func (n *Node) heartbeatLoop() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		n.m.RLock()
		peers := make([]*Peer, 0, len(n.Peers))
		for _, p := range n.Peers {
			peers = append(peers, p)
		}
		n.m.RUnlock()

		for _, p := range peers {
			url := fmt.Sprintf("%s/heartbeat", p.Address)
			resp, err := n.Client.Post(url, "application/json", nil)
			if err != nil || resp.StatusCode != http.StatusOK {
				log.Printf("[%s] Peer %s heartbeat failed", n.ID, p.ID)
				log.Println("If error :", err.Error())
				n.m.Lock()
				//  for now adding 10s as buffer time to see if node is active
				if time.Since(p.LastSeen) > 10*time.Second {
					// delete from peer map
					delete(n.Peers, p.ID)
					log.Printf("[%s] Peer %s removed", n.ID, p.ID)
				}
				n.m.Unlock()
				continue
			}
			log.Printf("[%s] Peer %s node active", n.ID, p.ID)
			resp.Body.Close()

			n.m.Lock()
			p.LastSeen = time.Now()
			p.Alive = true
			n.m.Unlock()
		}
	}
}
