// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(false)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

//TestServers list all servers by port
var TestServers = map[int]*TestServer{}

//TestServer is responsible for emulating a deepstream server
type TestServer struct {
	Port                 int
	Listener             net.Listener
	Server               http.Server
	Closer               io.Closer
	WebsocketConnections map[string]*websocket.Conn
	ActiveConnections    int
	upgrader             websocket.Upgrader
	ReceivedMessages     []string
	ReceivedErrors       []error
	shouldStop           bool
}

//Start the test server
func (ts *TestServer) Start() error {
	ts.shouldStop = false
	ts.WebsocketConnections = map[string]*websocket.Conn{}
	ts.upgrader = websocket.Upgrader{} // use default options
	fmt.Println("===========================> LISTENING <=======================", ts.Port)

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/deepstream", ts.handler)

	var err error
	ts.Closer, err = listenAndServeWithClose(fmt.Sprintf("127.0.0.1:%d", ts.Port), serverMux)
	if err != nil {
		return err
	}

	return nil
}

func listenAndServeWithClose(addr string, handler http.Handler) (io.Closer, error) {
	var (
		listener  net.Listener
		srvCloser io.Closer
		err       error
	)

	srv := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			log.Println("HTTP Server Error - ", err)
		}
	}()

	srvCloser = listener
	return srvCloser, nil
}

//Stop the test server
func (ts *TestServer) Stop() {
	//ts.shouldStop = true
	for _, ws := range ts.WebsocketConnections {
		ws.Close()
	}
	ts.WebsocketConnections = map[string]*websocket.Conn{}
	ts.ActiveConnections = 0
	if ts.Closer != nil {
		ts.Closer.Close()
		time.Sleep(500 * time.Millisecond) // Wait for ports to be reclaimed by OS
	}
	//ts.Listener.Close()
	//ts.Listener = nil
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
		c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
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
		if ts.shouldStop {
			return
		}
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
