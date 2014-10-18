package catcher

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

type Catcher struct {
	router    *mux.Router
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]*client
	broadcast chan *CaughtRequest
	logger    *logging.Logger
}

func NewCatcher() *Catcher {
	catcher := &Catcher{
		router: mux.NewRouter(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:   make(map[*websocket.Conn]*client),
		broadcast: make(chan *CaughtRequest),
		logger:    logging.MustGetLogger("request-catcher"),
	}
	catcher.init()
	return catcher
}

func (c *Catcher) init() {
	c.router.HandleFunc("/", c.indexHandler)
	c.router.HandleFunc("/init-client", c.initClient)
	c.router.PathPrefix("/static").Handler(http.FileServer(http.Dir("catcher/")))
	c.router.NotFoundHandler = http.HandlerFunc(c.catchRequests)
}

func (c *Catcher) Start() {
	go c.broadcaster()
	http.Handle("/", c.router)
	http.ListenAndServe(":4000", nil)
}

func (c *Catcher) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "catcher/templates/index.html")
}

func (c *Catcher) catchRequests(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found")
		return
	}
	c.broadcast <- convertRequest(r)
}

func (c *Catcher) broadcaster() {
	for req := range c.broadcast {
		for _, client := range c.clients {
			client.output <- req
		}
	}
}

func (c *Catcher) initClient(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Initializing a new client from %v", r.RemoteAddr)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	client := newClient(c, ws)
	c.clients[ws] = client
}
