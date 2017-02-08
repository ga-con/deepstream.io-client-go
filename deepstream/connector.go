// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
)

//OnMessageHandler represents a function that takes a message and does something
type OnMessageHandler func(*Message)

//Connector is an abstraction to the web socket connection to deepstream
type Connector struct {
	URL                 string
	ConnectionTimeoutMs int
	WriteTimeoutMs      int
	ReadTimeoutMs       int
	ConnectionLock      *sync.Mutex
	ConnectionState     interfaces.ConnectionState
	Client              *websocket.Conn
	MessageHandlers     []OnMessageHandler
	CloseHandler        func(string) error
}

//NewConnector creates a new connector to the specified Deepstream server url
func NewConnector(url string, connectionTimeoutMs, writeTimeoutMs, readTimeoutMs int, closeHandler func(string) error) *Connector {
	return &Connector{
		ConnectionTimeoutMs: connectionTimeoutMs,
		WriteTimeoutMs:      writeTimeoutMs,
		ReadTimeoutMs:       readTimeoutMs,
		ConnectionLock:      &sync.Mutex{},
		MessageHandlers:     []OnMessageHandler{},
		Client:              nil,
		ConnectionState:     interfaces.ConnectionStateAwaitingConnection,
		URL:                 url,
		CloseHandler:        closeHandler,
	}
}

//Connect to deepstream using websocket and starts monitoring traffic in the background
func (c *Connector) Connect() error {
	if c.ConnectionState != interfaces.ConnectionStateReconnecting {
		c.ConnectionState = interfaces.ConnectionStateAwaitingConnection
	}
	url := fmt.Sprintf("ws://%s/deepstream", c.URL)

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Duration(c.ConnectionTimeoutMs) * time.Millisecond
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.Client = client

	go func() {
		for {
			if c.Client == nil {
				return
			}

			messageType, msgBytes, err := c.Client.ReadMessage()
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				c.Close()
				c.ConnectionState = interfaces.ConnectionStateClosed
				if c.CloseHandler != nil {
					c.CloseHandler(string(err.Error()))
				}
				return
			}
			if err != nil {
				return
			}

			if messageType == websocket.BinaryMessage {
				//ON ERROR?
				//return fmt.Errorf("Message not understood")
				return
			}

			func() {
				messages, err := ParseMessages(string(msgBytes))
				if err != nil {
					//onErr?
					return
				}
				for _, msg := range messages {
					c.OnMessage(msg)
				}
			}()
		}
	}()

	return nil
}

func (c *Connector) acquireLock() func() {
	c.ConnectionLock.Lock()
	return func() {
		c.ConnectionLock.Unlock()
	}
}

//AddMessageHandler to handle incoming messages
func (c *Connector) AddMessageHandler(handler OnMessageHandler) {
	defer c.acquireLock()
	c.MessageHandlers = append(c.MessageHandlers, handler)
}

//OnMessage ensures all handlers are called
func (c *Connector) OnMessage(msg *Message) {
	for _, handler := range c.MessageHandlers {
		handler(msg)
	}
}

//WriteMessage sends a message to deepstream websocket connection using text
func (c *Connector) WriteMessage(msg []byte, msgTypeOrNil ...int) error {
	c.Client.SetWriteDeadline(time.Now().Add(300 * time.Millisecond))
	msgType := websocket.TextMessage
	if len(msgTypeOrNil) == 1 {
		msgType = msgTypeOrNil[1]
	}
	return c.Client.WriteMessage(msgType, msg)
}

//Close the connection to deepstream
func (c *Connector) Close() error {
	if c.Client != nil {
		err := c.Client.Close()
		if err != nil {
			return err
		}
		c.Client = nil
	}
	c.ConnectionState = interfaces.ConnectionStateClosed
	return nil
}
