// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import (
	"encoding/json"
	"fmt"

	"github.com/heynemann/deepstream.io-client-go/interfaces"
	"github.com/heynemann/deepstream.io-client-go/message"
)

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
		var err error
		proto, err = NewWebsocketProtocol(url)
		if err != nil {
			return nil, err
		}
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
	c.ConnectionState = interfaces.ConnectionStateAwaitingConnection

	err := c.Protocol.Connect()
	if err != nil {
		return c.Error(err)
	}

	//Send Challenge Response
	err = c.sendChallengeResponse()
	if err != nil {
		return c.Error(err)
	}

	err = c.receiveAck("C")
	if err != nil {
		return c.Error(err)
	}

	return nil
}

func (c *Client) sendChallengeResponse() error {
	challenge := message.NewChallengeResponseAction(c.URL)
	return c.Protocol.SendAction(challenge)
}

func (c *Client) receiveAck(expectedTopic string) error {
	// Receive connection Ack
	actions, err := c.Protocol.RecvActions()
	if err != nil {
		return err
	}

	if len(actions) != 1 {
		//TODO: change this
		return fmt.Errorf("Expected Ack")
	}

	action := actions[0]
	if a, ok := action.(*message.AckAction); !ok || a.Topic != expectedTopic {
		return fmt.Errorf("Expected Ack with topic %s.", expectedTopic)
	}

	return nil
}

//Close connection to deepstream.io server
func (c *Client) Close() error {
	err := c.Protocol.Close()
	if err != nil {
		return c.Error(err)
	}

	c.ConnectionState = interfaces.ConnectionStateClosed
	return nil
}

//Login with deepstream.io server
func (c *Client) Login(authParams map[string]interface{}) error {
	params, err := json.Marshal(authParams)
	if err != nil {
		return err
	}
	msg := &message.Message{
		Topic:   interfaces.TopicAuth,
		Action:  interfaces.ActionRequest,
		RawData: []string{string(params)},
	}
	authRequestAction, err := message.NewAuthRequestAction(msg)
	if err != nil {
		return err
	}

	//Send Authentication Request
	err = c.Protocol.SendAction(authRequestAction)
	if err != nil {
		return err
	}

	// Receive authentication Ack
	actions, err := c.Protocol.RecvActions()
	if len(actions) != 1 {
		//TODO: change this
		return c.Error(fmt.Errorf("Expected Auth Ack"))
	}

	action := actions[0]
	if a, ok := action.(*message.AckAction); !ok || a.Topic != "A" {
		return c.Error(fmt.Errorf("Expected Auth Ack"))
	}

	c.ConnectionState = interfaces.ConnectionStateOpen

	return nil
}

//Error handlers errors in client
func (c *Client) Error(err error) error {
	c.ConnectionState = interfaces.ConnectionStateError
	return err
}
