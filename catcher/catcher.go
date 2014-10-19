package catcher

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

type Catcher struct {
	host     string
	port     int
	router   *mux.Router
	upgrader websocket.Upgrader
	hosts    map[string]*Host
	logger   *logging.Logger
}

func NewCatcher(host string, port int) *Catcher {
	catcher := &Catcher{
		host:   host,
		port:   port,
		router: mux.NewRouter(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hosts:  make(map[string]*Host),
		logger: logging.MustGetLogger("request-catcher"),
	}
	catcher.init()
	return catcher
}

func (c *Catcher) init() {
	c.router.HandleFunc("/", c.rootHandler).Host(c.host)
	c.router.HandleFunc("/", c.indexHandler)
	c.router.HandleFunc("/init-client", c.initClient)
	c.router.PathPrefix("/static").Handler(http.FileServer(http.Dir("catcher/")))
	c.router.NotFoundHandler = http.HandlerFunc(c.catchRequests)
}

func (c *Catcher) Start() {
	http.Handle("/", c.router)
	fullHost := c.host + ":" + strconv.Itoa(c.port)
	c.logger.Info("Listening on %v on port %v", c.host, c.port)
	http.ListenAndServe(fullHost, nil)
}

func (c *Catcher) getHost(hostString string) *Host {
	hostString = hostWithoutPort(hostString)

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
	http.ServeFile(w, r, "catcher/templates/index.html")
}

func (c *Catcher) catchRequests(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found")
		return
	}
	caughtRequest := convertRequest(r)
	host := c.getHost(caughtRequest.Host)
	c.logger.Info("Routing caught request to %v", host)
	host.broadcast <- caughtRequest
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

	clientHost := c.getHost(r.Host)
	c.logger.Info("Initializing a new client on host %v", clientHost.Host)
	clientHost.addClient(newClient(c, clientHost, ws))
}
