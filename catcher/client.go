package catcher

import (
	"time"

	"github.com/gorilla/websocket"
)

const outputChannelBuffer = 5
const pingFrequency = 10 * time.Second
const writeWait = 5 * time.Second
const pongWait = 60 * time.Second
const maxMessageSize = 1024

type client struct {
	pingTicker *time.Ticker
	catcher    *Catcher
	host       *Host
	conn       *websocket.Conn
	closed     chan struct{}
	output     chan interface{}
}

func newClient(catcher *Catcher, host *Host, conn *websocket.Conn) *client {
	c := &client{
		pingTicker: time.NewTicker(pingFrequency),
		catcher:    catcher,
		host:       host,
		conn:       conn,
		closed:     make(chan struct{}),
		output:     make(chan interface{}, outputChannelBuffer),
	}
	go c.writeLoop()
	go c.readLoop()
	return c
}

func (c *client) ping() error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (c *client) Close() error {
	select {
	case <-c.closed:
		// already closed
		return nil
	default:
	}

	c.host.clients.Delete(c.conn)
	close(c.closed)
	c.pingTicker.Stop()

	// Be nice and issue a CloseMessage before closing the conn.
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})

	c.conn.Close()
	return err
}

func (c *client) sendJSON(obj interface{}) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteJSON(obj)
}

func (c *client) writeLoop() {
	defer c.Close()

	for {
		select {
		case <-c.pingTicker.C:
			if err := c.ping(); err != nil {
				c.catcher.logger.Errorf("Error pinging: %v", err)
				return
			}
		case msg := <-c.output:
			if err := c.sendJSON(msg); err != nil {
				c.catcher.logger.Errorf("Error sending message: %v", err)
				return
			}
		case <-c.closed:
			return
		}
	}
}

func (c *client) readLoop() {
	defer c.Close()

	// We don't care about what the client sends to us, but we need to
	// read it to keep the connection fresh.
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(msg string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			break
		}
	}
}
