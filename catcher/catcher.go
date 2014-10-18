package catcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Catcher struct {
	router    *mux.Router
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]*client
	broadcast chan *CaughtRequest
}

type CaughtRequest struct {
	Time          time.Time   `json:"time"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Headers       http.Header `json:"headers"`
	ContentLength int64       `json:"content_length"`
	RemoteAddr    string      `json:"remote_addr"`
	Form          url.Values  `json:"form_values"`
	Body          string      `json:"body"`
}

type client struct {
	conn   *websocket.Conn
	output chan interface{}
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
	fmt.Printf("Initializing a new client.\n")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &client{
		conn:   ws,
		output: make(chan interface{}, 1),
	}
	c.clients[ws] = client

	go c.writePump(client)
	// We don't care about what the client sends to us, but we need to
	// read it to keep the connection fresh.
	go func() {
		ws.SetReadLimit(1024)
		ws.SetReadDeadline(time.Now().Add(time.Minute))
		ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(time.Minute)); return nil })
		for {
			if _, _, err := ws.NextReader(); err != nil {
				ws.Close()
				break
			}
		}
	}()
}

func (catcher *Catcher) writePump(client *client) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	defer client.conn.Close()
	defer delete(catcher.clients, client.conn)
	for {
		select {
		case <-ticker.C:
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
		case msg, ok := <-client.output:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.conn.WriteJSON(msg); err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
		}
	}
}

func convertRequest(req *http.Request) *CaughtRequest {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
	}

	r := &CaughtRequest{
		Time:          time.Now(),
		Method:        req.Method,
		Path:          req.RequestURI,
		Headers:       req.Header,
		ContentLength: req.ContentLength,
		RemoteAddr:    req.RemoteAddr,
		Form:          req.PostForm,
		Body:          string(body),
	}
	return r
}
