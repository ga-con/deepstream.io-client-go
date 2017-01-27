// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import (
	"fmt"
	"time"

	"github.com/heynemann/deepstream.io-client-go/interfaces"
	"github.com/heynemann/deepstream.io-client-go/message"
	"github.com/looplab/fsm"
)

//Client represents a connection to a deepstream.io server
type Client struct {
	URL             string
	Protocol        interfaces.Protocol
	ConnectionState interfaces.ConnectionState
	FSM             *fsm.FSM
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

	cli.configureFSM()

	err := cli.Connect()
	if err != nil {
		return cli, err
	}

	return cli, nil
}

func (c *Client) configureFSM() error {
	c.FSM = fsm.NewFSM(
		string(interfaces.ConnectionStateAwaitingConnection),
		fsm.Events{
			{
				Name: "connect", Src: []string{
					string(interfaces.ConnectionStateAwaitingConnection),
				}, Dst: string(interfaces.ConnectionStateConnecting),
			},
			{
				Name: "connected", Src: []string{
					string(interfaces.ConnectionStateConnecting),
				},
				Dst: string(interfaces.ConnectionStateConnected),
			},
			{
				Name: "challengeReceived", Src: []string{
					string(interfaces.ConnectionStateConnected),
				},
				Dst: string(interfaces.ConnectionStateChallengeReceived),
			},
			{
				Name: "challenge", Src: []string{
					string(interfaces.ConnectionStateChallengeReceived),
				},
				Dst: string(interfaces.ConnectionStateChallenging),
			},
			{
				Name: "authenticationRequested", Src: []string{
					string(interfaces.ConnectionStateChallenging),
				},
				Dst: string(interfaces.ConnectionStateAwaitingAuthentication),
			},
			{
				Name: "authenticate", Src: []string{
					string(interfaces.ConnectionStateAwaitingAuthentication),
				},
				Dst: string(interfaces.ConnectionStateOpen),
			},
			{
				Name: "close", Src: []string{
					string(interfaces.ConnectionStateOpen), string(interfaces.ConnectionStateError),
				},
				Dst: string(interfaces.ConnectionStateClosed),
			},
		},
		fsm.Callbacks{
			"enter_state": c.onStateChange,
		},
	)

	return nil
}

func (c *Client) onStateChange(e *fsm.Event) {
	fmt.Println(e.Src, e.Dst)
	var err error
	switch e.Dst {
	case string(interfaces.ConnectionStateConnecting):
		fmt.Println("Connecting")
		err = c.handleConnecting()
	case string(interfaces.ConnectionStateConnected):
		fmt.Println("Connected")
		err = c.handleConnected()
	case string(interfaces.ConnectionStateChallengeReceived):
		fmt.Println("Challenge Received")
		err = c.handleChallengeReceived()
	case string(interfaces.ConnectionStateChallenging):
		fmt.Println("Challenging")
		err = c.handleChallenging()
	}
	e.Async()

	if err != nil {
		//TODO: do what?
	}
}

func (c *Client) handleConnecting() error {
	err := c.Protocol.Connect()
	if err != nil {
		return c.Error(err)
	}

	c.FSM.Event("connected")
	return nil
}

func (c *Client) handleConnected() error {
	actions, err := c.Protocol.RecvActions()
	if err != nil {
		return err
	}
	if len(actions) != 1 {
		//TODO: Change this
		return fmt.Errorf("authentication challenge expected")
	}
	action := actions[0]
	if _, ok := action.(*message.ChallengeAction); !ok {
		return fmt.Errorf("authentication challenge expected 2")
	}

	c.FSM.Event("challengeReceived")

	return nil
}

func (c *Client) handleChallengeReceived() error {
	err := c.sendChallengeResponse()
	if err != nil {
		return err
	}

	c.FSM.Event("challenging")

	return nil
}

func (c *Client) handleChallenging() error {
	err := c.receiveAck("C")
	if err != nil {
		return err
	}
	//c.FSM.Event("challengeReceived")

	return nil
}

//Connect with deepstream.io server
func (c *Client) Connect() error {
	err := c.FSM.Event("connect")
	if err != nil {
		return err
	}

	for c.FSM.Current() != string(interfaces.ConnectionStateAwaitingAuthentication) {
		time.Sleep(1 * time.Millisecond)
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
	//params, err := json.Marshal(authParams)
	//if err != nil {
	//return err
	//}
	//msg := &message.Message{
	//Topic:   interfaces.TopicAuth,
	//Action:  interfaces.ActionRequest,
	//RawData: []string{string(params)},
	//}
	//authRequestAction, err := message.NewAuthRequestAction(msg)
	//if err != nil {
	//return err
	//}

	////Send Authentication Request
	//err = c.Protocol.SendAction(authRequestAction)
	//if err != nil {
	//return err
	//}

	//// Receive authentication Ack
	//actions, err := c.Protocol.RecvActions()
	//if len(actions) != 1 {
	////TODO: change this
	//return c.Error(fmt.Errorf("Expected Auth Ack"))
	//}

	//action := actions[0]
	//if a, ok := action.(*message.AckAction); !ok || a.Topic != "A" {
	//return c.Error(fmt.Errorf("Expected Auth Ack"))
	//}

	//c.ConnectionState = interfaces.ConnectionStateOpen

	return nil
}

//Error handlers errors in client
func (c *Client) Error(err error) error {
	//c.ConnectionState = interfaces.ConnectionStateError
	return err
}
