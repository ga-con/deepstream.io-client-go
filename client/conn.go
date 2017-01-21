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
	URL             string
	Protocol        interfaces.Protocol
	ConnectionState interfaces.ConnectionState
}

//New creates a new client
func New(url string, protocolOrNil ...interfaces.Protocol) (*Client, error) {
	var proto interfaces.Protocol
	if len(protocolOrNil) == 1 {
		proto = protocolOrNil[0]
	}
	if proto == nil {
		//TODO: load real protocol
	}
	cli := &Client{
		URL:             url,
		ConnectionState: interfaces.ConnectionStateClosed,
		Protocol:        proto,
	}

	err := cli.Connect()
	if err != nil {
		return cli, err
	}

	return cli, nil
}

//Connect with deepstream.io server
func (c *Client) Connect() error {
	err := c.Protocol.Connect()
	if err != nil {
		return c.Error(err)
	}

	c.ConnectionState = interfaces.ConnectionStateAwaitingAuthentication
	return nil
}

//Login with deepstream.io server
func (c *Client) Login(authParams map[string]interface{}) error {
	return c.Protocol.Authenticate(authParams)
}

//Error handlers errors in client
func (c *Client) Error(err error) error {
	c.ConnectionState = interfaces.ConnectionStateError
	return err
}
