package catcher

import "github.com/gorilla/websocket"

// Host represents a host on which we've received requests.
type Host struct {
	Host      string
	clients   map[*websocket.Conn]*client
	broadcast chan *CaughtRequest
}

func newHost(host string) *Host {
	hostObj := &Host{
		Host:      host,
		clients:   make(map[*websocket.Conn]*client),
		broadcast: make(chan *CaughtRequest),
	}
	go hostObj.broadcaster()
	return hostObj
}

func (h *Host) broadcaster() {
	for req := range h.broadcast {
		for _, client := range h.clients {
			client.output <- req
		}
	}
}

func (h *Host) addClient(c *client) {
	h.clients[c.conn] = c
}
