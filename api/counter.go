package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// IncrementHandler, Increment local counter
func (n *Node) IncrementHandler(w http.ResponseWriter, r *http.Request) {
	n.m.Lock()
	n.Counter++
	val := n.Counter
	n.m.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"value": val,
	})

	// Sync with peers
	go n.BroadcastUpdate()
}

// DecrementHandler, Decrement local counter
func (n *Node) DecrementHandler(w http.ResponseWriter, r *http.Request) {
	n.m.Lock()
	n.Counter--
	val := n.Counter
	n.m.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"value": val,
	})

	// Sync with peers
	go n.BroadcastUpdate()
}

// ValueHandler, get local counter count
func (n *Node) ValueHandler(w http.ResponseWriter, r *http.Request) {
	n.m.RLock()
	val := n.Counter
	n.m.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"value": val,
	})
}

// Payload from peers
type CounterUpdate struct {
	Value int `json:"value"`
}

func (n *Node) SyncHandler(w http.ResponseWriter, r *http.Request) {
	var update CounterUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n.m.Lock()
	n.Counter = update.Value
	n.m.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (n *Node) BroadcastUpdate() {
	n.m.RLock()
	val := n.Counter
	n.m.RUnlock()

	update := CounterUpdate{Value: val}
	body, _ := json.Marshal(update)

	for _, peer := range n.Peers {
		go func(p *Peer) {
			// calling sync API for peer services
			resp, err := n.Client.Post(p.Address+"/counter/sync", "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Printf("failed to sync with %s: %v", p.ID, err)
				return
			}
			resp.Body.Close()
		}(peer)
	}
}
