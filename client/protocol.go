// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

//WebsocketProtocol is the default protocol for deepstream.io
type WebsocketProtocol struct {
	Host   string
	Path   string
	Client *websocket.Conn
}

//NewWebsocketProtocol creates a new instance
func NewWebsocketProtocol(host, path string) (*WebsocketProtocol, error) {
	ws := &WebsocketProtocol{
		Host: host,
		Path: path,
	}
	return ws, nil
}

//Connect to deepstream.io
func (w *WebsocketProtocol) Connect() error {
	u := url.URL{Scheme: "ws", Host: w.Host, Path: w.Path}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	w.Client = c
	return nil
}

//Close websocket connection
func (w *WebsocketProtocol) Close() error {
	if w.Client == nil {
		return nil
	}

	err := w.Client.Close()
	if err != nil {
		return err
	}

	w.Client = nil
	return err
}
