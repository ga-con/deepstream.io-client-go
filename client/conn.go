// deepstream.io-client-go
// https://github.com/ga-con/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/ga-con/deepstream.io-client-go/interfaces"
	"github.com/ga-con/deepstream.io-client-go/message"
)

//Client represents a connection to a deepstream.io server
type Client struct {
	URL string
	//Protocol        interfaces.Protocol
	Conn            *websocket.Conn
	ConnectionState interfaces.ConnectionState
}

//New creates a new client
func New(url string, protocolOrNil ...interfaces.Protocol) (*Client, error) {
	/*var proto interfaces.Protocol
	if len(protocolOrNil) == 1 {
		proto = protocolOrNil[0]
	}*/
	/*if proto == nil {
		var err error
		proto, err = NewWebsocketProtocol(url)
		if err != nil {
			return nil, err
		}
	}*/
	cli := &Client{
		URL:             url,
		ConnectionState: interfaces.ConnectionStateClosed,
		Conn:            &websocket.Conn{},
		//Protocol:        proto,
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

	url := fmt.Sprintf("ws://%s/deepstream", c.URL)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.Conn = conn

	err = c.getAuthChallenge()
	if err != nil {
		return err
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
	return c.SendAction(challenge)
}

func (c *Client) receiveAck(expectedTopic string) error {
	// Receive connection Ack
	actions, err := c.RecvActions()
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
	fmt.Println("------------CLOSED")

	if c.Conn == nil {
		return nil
	}

	err := c.Conn.Close()
	if err != nil {
		return err
	}

	c.Conn = nil

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
	err = c.SendAction(authRequestAction)
	if err != nil {
		return err
	}

	// Receive authentication Ack
	actions, err := c.RecvActions()
	if len(actions) != 1 {
		//TODO: change this
		return c.Error(fmt.Errorf("Expected Auth Ack"))
	}

	action := actions[0]
	if a, ok := action.(*message.AckAction); !ok || a.Topic != "A" {
		return c.Error(fmt.Errorf("Expected Auth Ack"))
	}

	c.ConnectionState = interfaces.ConnectionStateOpen

	// Listen RecvActions
	go func() {
		for {
			acts, err := c.RecvActions()
			fmt.Println("----------c.Protocol.RecvActions:")

			if err != nil {
				fmt.Errorf("handlerConnection:%v", err)
				return
			}

			for _, act := range acts {
				fmt.Println("act:", act)

				if a, ok := act.(*message.PingAction); !ok || a.Topic != "C" {
					fmt.Print("handlerConnection:", act)
					continue
				}

				rAction, _ := message.NewPongAction(&message.Message{
					Topic:  interfaces.TopicConnection,
					Action: interfaces.ActionPong,
				})

				if err := c.SendAction(rAction); err != nil {
					fmt.Println("NewPongAction Error:", err)
				}
			}
		}
	}()

	return nil
}

//Error handlers errors in client
func (c *Client) Error(err error) error {
	c.ConnectionState = interfaces.ConnectionStateError
	return err
}

func (c *Client) getAuthChallenge() error {
	_, body, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}
	msgs, err := message.ParseMessages(string(body))
	if len(msgs) != 1 {
		//TODO: Change this
		return fmt.Errorf("authentication challenge expected")
	}
	msg := msgs[0]
	action, err := message.CathegorizeAction(msg)
	if err != nil {
		return err
	}
	if _, ok := action.(*message.ChallengeAction); !ok {
		return fmt.Errorf("authentication challenge expected 2")
	}

	return nil
}

//SendAction writes an action in the websocket stream
func (c *Client) SendAction(action interfaces.Action) error {
	msg := action.ToAction()
	err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

//RecvActions receives actions from the websocket stream
func (c *Client) RecvActions() ([]interfaces.Action, error) {
	_, body, err := c.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	msgs, err := message.ParseMessages(string(body))
	if err != nil {
		return nil, err
	}
	var actions []interfaces.Action
	for _, msg := range msgs {
		action, err := message.CathegorizeAction(msg)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	return actions, nil
}
