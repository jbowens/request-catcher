package catcher

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

type Catcher struct {
	config   *Configuration
	router   *mux.Router
	upgrader websocket.Upgrader
	logger   *logging.Logger

	hostsMu sync.Mutex
	hosts   map[string]*Host
}

func NewCatcher(config *Configuration) *Catcher {
	c := &Catcher{
		config: config,
		router: mux.NewRouter(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		logger: logging.MustGetLogger("request-catcher"),

		hosts: make(map[string]*Host),
	}
	c.router.HandleFunc("/", c.rootHandler).Host(c.config.RootHost)
	c.router.HandleFunc("/", c.indexHandler)
	c.router.HandleFunc("/init-client", c.initClient)
	c.router.PathPrefix("/static").Handler(http.FileServer(http.Dir("catcher/")))
	c.router.NotFoundHandler = http.HandlerFunc(c.catchRequests)
	return c
}

func (c *Catcher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.router.ServeHTTP(rw, req)
}

func (c *Catcher) host(hostString string) *Host {
	hostString = hostWithoutPort(hostString)

	c.hostsMu.Lock()
	defer c.hostsMu.Unlock()
	if host, ok := c.hosts[hostString]; ok {
		return host
	}
	host := newHost(hostString)
	c.hosts[hostString] = host
	return host
}

func (c *Catcher) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "catcher/templates/root.html")
}

func (c *Catcher) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Some people mistakenly expect requests to the index of the subdomain
	// to be caught. For now, just catch those as well. Later I should move
	// the index to be hosted at requestcatcher.com.
	c.catch(r)

	http.ServeFile(w, r, "catcher/templates/index.html")
}

func (c *Catcher) catchRequests(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found")
		return
	}

	c.catch(r)

	// Respond to the request
	fmt.Fprintf(w, "request caught")
}

func (c *Catcher) initClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	clientHost := c.host(r.Host)
	c.logger.Info("Initializing a new client on host %v", clientHost.Host)
	clientHost.clients.Store(c, newClient(c, clientHost, ws))
}

func (c *Catcher) catch(r *http.Request) {
	caughtRequest := convertRequest(r)

	// Broadcast it to everyone listening for requests on this host
	host := c.host(caughtRequest.Host)
	host.broadcast <- caughtRequest
}
