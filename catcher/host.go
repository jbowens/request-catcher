package catcher

import "sync"

// Host represents a host on which we've received requests.
type Host struct {
	Host      string
	broadcast chan *CaughtRequest

	// clients is a map from the pointer to the websocket
	// connection to the pointer to the corresponding client
	// struct. It's a sync.Map because sync.Map.Range doesn't
	// need to keep a mutex locked during iteration, which
	// is good for us if a client is being slow to respond.
	clients sync.Map // map[*websocket.Conn]*client
}

func newHost(host string) *Host {
	hostObj := &Host{
		Host:      host,
		broadcast: make(chan *CaughtRequest),
	}
	go hostObj.broadcaster()
	return hostObj
}

func (h *Host) broadcaster() {
	for req := range h.broadcast {
		h.clients.Range(func(conn, untypedClient interface{}) bool {
			typedClient := untypedClient.(*client)
			typedClient.output <- req
			return true
		})
	}
}
