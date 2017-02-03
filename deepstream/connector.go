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
}

//NewConnector creates a new connector to the specified Deepstream server url
func NewConnector(url string, connectionTimeoutMs, writeTimeoutMs, readTimeoutMs int) *Connector {
	return &Connector{
		ConnectionTimeoutMs: connectionTimeoutMs,
		WriteTimeoutMs:      writeTimeoutMs,
		ReadTimeoutMs:       readTimeoutMs,
		ConnectionLock:      &sync.Mutex{},
		MessageHandlers:     []OnMessageHandler{},
		Client:              nil,
		ConnectionState:     interfaces.ConnectionStateInitializing,
		URL:                 url,
	}
}

//Connect to deepstream using websocket and starts monitoring traffic in the background
func (c *Connector) Connect() error {
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
			//c.Client.SetReadDeadline(time.Now().Add(time.Duration(c.ReadTimeoutMs) * time.Millisecond))
			messageType, msgBytes, err := c.Client.ReadMessage()
			if err != nil {
				//ON ERROR?
				continue
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
