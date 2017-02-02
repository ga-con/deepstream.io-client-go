// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream

import (
	"fmt"

	"github.com/heynemann/deepstream.io-client-go/interfaces"
	"github.com/heynemann/deepstream.io-client-go/message"
)

//ClientOptions used to connect to deepstream
type ClientOptions struct {
	AutoReconnect bool
	AutoLogin     bool
	Username      string
	Password      string
}

//DefaultOptions to connect to deepstream
func DefaultOptions() *ClientOptions {
	return &ClientOptions{
		AutoReconnect: true,
		AutoLogin:     true,
		Username:      "",
		Password:      "",
	}
}

//Client represents a connection to a deepstream.io server
type Client struct {
	AuthParams     map[string]interface{}
	Connector      *Connector
	Options        *ClientOptions
	loginRequested bool
}

//New creates a new client
func New(url string, optionsOrNil ...*ClientOptions) (*Client, error) {
	options := DefaultOptions()
	if len(optionsOrNil) == 1 {
		options = optionsOrNil[0]
	}
	username := options.Username
	password := options.Password
	authParams := map[string]interface{}{}
	if username != "" {
		authParams = map[string]interface{}{
			"username": username,
			"password": password,
		}
	}
	cli := &Client{
		Connector:      NewConnector(url),
		Options:        options,
		AuthParams:     authParams,
		loginRequested: false,
	}

	cli.Connector.AddMessageHandler(cli.onMessage)

	err := cli.startMonitoringConnection()
	if err != nil {
		return cli, err
	}

	return cli, nil
}

func (c *Client) startMonitoringConnection() error {
	err := c.Connector.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) onMessage(msg *message.Message) {
	//fmt.Println(msg.Topic, msg.Action)
	var err error
	switch {
	case msg.Topic == "C":
		err = c.handleConnectionMessages(msg)
	case msg.Topic == "A":
		err = c.handleAuthenticationMessages(msg)
	}

	if err != nil {
		//On error?
	}
}

//GetConnectionState returns the connection state for the connector
func (c *Client) GetConnectionState() interfaces.ConnectionState {
	return c.Connector.ConnectionState
}

func (c *Client) handleConnectionMessages(msg *message.Message) error {
	switch {
	case msg.Action == "CH":
		return c.handleChallengeRequest(msg)
	case msg.Action == "A":
		if c.Connector.ConnectionState == interfaces.ConnectionStateChallenging {
			return c.handleChallengeAck(msg)
		}
	default:
		fmt.Println("Message not understood!")
	}

	return nil
}

func (c *Client) handleChallengeRequest(msg *message.Message) error {
	c.Connector.ConnectionState = interfaces.ConnectionStateChallenging
	challenge := message.NewChallengeResponseAction(c.Connector.URL)
	return c.Connector.WriteMessage([]byte(challenge.ToAction()))
}

func (c *Client) handleChallengeAck(msg *message.Message) error {
	c.Connector.ConnectionState = interfaces.ConnectionStateAwaitingConnection
	if c.Options.AutoLogin || c.loginRequested {
		return c.Login()
	}
	return nil
}

//Login to deepstream - if connection is still being started, it will login as soon as possible
func (c *Client) Login() error {
	state := c.GetConnectionState()
	if !c.Options.AutoLogin && (state == interfaces.ConnectionStateChallenging ||
		state == interfaces.ConnectionStateInitializing) {
		c.loginRequested = true
		return nil
	}

	if state != interfaces.ConnectionStateAwaitingConnection {
		return c.Error(fmt.Errorf("The connection should be restored before logging in (%s).", state))
	}

	c.loginRequested = false

	authRequestAction, err := message.NewAuthRequestAction(c.AuthParams)
	if err != nil {
		return err
	}

	c.Connector.ConnectionState = interfaces.ConnectionStateAuthenticating

	//Send Authentication Request
	return c.Connector.WriteMessage([]byte(authRequestAction.ToAction()))
}

func (c *Client) handleAuthenticationMessages(msg *message.Message) error {
	switch {
	case msg.Action == "A":
		if c.Connector.ConnectionState == interfaces.ConnectionStateAuthenticating {
			return c.handleAuthenticationAck(msg)
		}
	case msg.Action == "E":
		return c.Error(fmt.Errorf("Could not connect to deepstream.io server with the provided credentials (user: %s).", c.AuthParams["user"]))
	default:
		fmt.Println("Message not understood!")
	}

	return nil
}

func (c *Client) handleAuthenticationAck(msg *message.Message) error {
	c.Connector.ConnectionState = interfaces.ConnectionStateOpen
	return nil
}

//Error handlers errors in client
func (c *Client) Error(err error) error {
	c.Connector.ConnectionState = interfaces.ConnectionStateError
	return err
}
