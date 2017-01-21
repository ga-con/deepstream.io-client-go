// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import "github.com/heynemann/deepstream.io-client-go/interfaces"

//Client represents a connection to a deepstream.io server
type Client struct {
	URL      string
	Protocol interfaces.Protocol
}

//New creates a new client
func New(url string, protocolOrNil ...interfaces.Protocol) *Client {
	var proto interfaces.Protocol
	if len(protocolOrNil) == 1 {
		proto = protocolOrNil[0]
	}
	if proto == nil {
		//TODO: load real protocol
	}
	return &Client{
		URL:      url,
		Protocol: proto,
	}
}

//Login with deepstream.io server
func (c *Client) Login() error {
	return c.Protocol.Authenticate()
}
