package catcher

import (
	"fmt"
	"net/http"
	"strings"
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
	c.router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("catcher/static"))))
	c.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "catcher/static/favicon.ico")
	})
	c.router.NotFoundHandler = http.HandlerFunc(c.catchRequests)
	return c
}

func (c *Catcher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.Host, "www.") {
		rw.Header().Set("Connection", "close")
		url := "https://" + strings.TrimPrefix(req.Host, "www.") + req.URL.String()
		http.Redirect(rw, req, url, http.StatusMovedPermanently)
		return
	}

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
	c.Catch(r)

	http.ServeFile(w, r, "catcher/templates/index.html")
}

func (c *Catcher) catchRequests(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found")
		return
	}

	c.Catch(r)

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
	c.logger.Infof("Initializing a new client on host %v", clientHost.Host)
	clientHost.clients.Store(c, newClient(c, clientHost, ws))
}

func (c *Catcher) Catch(r *http.Request) {
	hostString := hostWithoutPort(r.Host)
	c.hostsMu.Lock()
	host, ok := c.hosts[hostString]
	c.hostsMu.Unlock()

	if !ok {
		// No one is listening, so no reason to catch it.
		return
	}

	// Broadcast it to everyone listening for requests on this host
	caughtRequest := convertRequest(r)
	host.broadcast <- caughtRequest
}
