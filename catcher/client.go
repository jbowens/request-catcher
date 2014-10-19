package catcher

import (
	"time"

	"github.com/gorilla/websocket"
)

const outputChannelBuffer = 5
const pingFrequency = 10 * time.Second
const writeWait = 5 * time.Second
const maxMessageSize = 1024

type client struct {
	pingTicker *time.Ticker
	catcher    *Catcher
	host       *Host
	conn       *websocket.Conn
	output     chan interface{}
}

func newClient(catcher *Catcher, host *Host, conn *websocket.Conn) *client {
	c := &client{
		pingTicker: time.NewTicker(pingFrequency),
		catcher:    catcher,
		host:       host,
		conn:       conn,
		output:     make(chan interface{}, outputChannelBuffer),
	}
	go c.writeLoop()
	go c.readLoop()
	return c
}

func (c *client) ping() error {
	c.catcher.logger.Info("Pinging a client")
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (c *client) Close() error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *client) sendJSON(obj interface{}) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteJSON(obj)
}

func (c *client) writeLoop() {
	defer func() {
		c.catcher.logger.Info("Client exiting")
		c.pingTicker.Stop()
		c.conn.Close()
		delete(c.host.clients, c.conn)
	}()

	for {
		select {
		case <-c.pingTicker.C:
			if err := c.ping(); err != nil {
				c.catcher.logger.Error("Error pinging: %v", err)
				return
			}
		case msg, ok := <-c.output:
			if !ok {
				c.Close()
				return
			}

			if err := c.sendJSON(msg); err != nil {
				c.catcher.logger.Error("Error sending message: %v", err)
				return
			}
		}
	}
}

func (c *client) readLoop() {
	// We don't care about what the client sends to us, but we need to
	// read it to keep the connection fresh.
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Time{})
	c.conn.SetPongHandler(func(msg string) error {
		c.catcher.logger.Debug("Pong from a client")
		return nil
	})
	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			c.conn.Close()
			break
		}
	}
}
