// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

//TestServers list all servers by port
var TestServers = map[int]*TestServer{}

//TestServer is responsible for emulating a deepstream server
type TestServer struct {
	Port                 int
	Listener             net.Listener
	WebsocketConnections map[string]*websocket.Conn
	ActiveConnections    int
	upgrader             websocket.Upgrader
	ReceivedMessages     []string
	ReceivedErrors       []error
}

//Start the test server
func (ts *TestServer) Start() error {
	ts.WebsocketConnections = map[string]*websocket.Conn{}
	ts.upgrader = websocket.Upgrader{} // use default options
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", ts.Port))
	if err != nil {
		return err
	}
	ts.Listener = ln

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/deepstream", ts.handler)

	server := http.Server{
		Handler: serverMux,
	}

	go func() {
		server.Serve(ln)
	}()

	return nil
}

//Stop the test server
func (ts *TestServer) Stop() {
	ts.Listener.Close()
}

func (ts *TestServer) handler(w http.ResponseWriter, r *http.Request) {
	ts.ActiveConnections++
	c, err := ts.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade:", err)
		return
	}
	id := uuid.NewV4().String()
	ts.WebsocketConnections[id] = c

	for {
		_, msg, err := c.ReadMessage()

		if websocket.IsCloseError(err) {
			c.Close()
			delete(ts.WebsocketConnections, id)
			ts.ActiveConnections--
			return
		}
		if err != nil {
			delete(ts.WebsocketConnections, id)
			ts.ActiveConnections--
			return
		}

		ts.ReceivedMessages = append(ts.ReceivedMessages, string(msg))
	}
}

//SendMessage to all registered servers
func (ts *TestServer) SendMessage(message string) error {
	var err error
	for _, conn := range ts.WebsocketConnections {
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

//HasMessage indicates whether the server received the specified message
func (ts *TestServer) HasMessage(expectedMessage string) error {
	for _, msg := range ts.ReceivedMessages {
		if msg == expectedMessage {
			return nil
		}
	}

	return fmt.Errorf("Message '%s' was not received by the server.", expectedMessage)
}

//StartTestServer at given port
func StartTestServer(port int) error {
	if _, ok := TestServers[port]; ok {
		return nil
	}

	ts := &TestServer{
		Port:             port,
		ReceivedMessages: []string{},
		ReceivedErrors:   []error{},
	}
	err := ts.Start()
	if err != nil {
		return err
	}

	TestServers[port] = ts
	return nil
}

//ResetTestServers resets everything
func ResetTestServers() {
	TestServers = map[int]*TestServer{}
}
