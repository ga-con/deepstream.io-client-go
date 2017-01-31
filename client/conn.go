// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/heynemann/deepstream.io-client-go/interfaces"
	"github.com/heynemann/deepstream.io-client-go/message"
)

type OnMessageHandler func(*message.Message)

type Connector struct {
	URL             string
	ConnectionLock  *sync.Mutex
	ConnectionState interfaces.ConnectionState
	Client          *websocket.Conn
	MessageHandlers []OnMessageHandler
}

func NewConnector(url string) *Connector {
	return &Connector{
		ConnectionLock:  &sync.Mutex{},
		MessageHandlers: []OnMessageHandler{},
		Client:          nil,
		ConnectionState: interfaces.ConnectionStateClosed,
		URL:             url,
	}
}

func (c *Connector) Connect() error {
	url := fmt.Sprintf("ws://%s/deepstream", c.URL)
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.Client = client

	go func() {
		for {
			messageType, msgBytes, err := c.Client.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.BinaryMessage {
				//ON ERROR?
				//return fmt.Errorf("Message not understood")
				return
			}

			func() {
				messages, err := message.ParseMessages(string(msgBytes))
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

func (c *Connector) AddMessageHandler(handler OnMessageHandler) {
	defer c.acquireLock()
	c.MessageHandlers = append(c.MessageHandlers, handler)
}

func (c *Connector) OnMessage(msg *message.Message) {
	for _, handler := range c.MessageHandlers {
		handler(msg)
	}
}

func (c *Connector) WriteMessage(msg []byte) error {
	return c.Client.WriteMessage(websocket.TextMessage, msg)
}

type ClientOptions struct {
	AutoReconnect bool
	AutoLogin     bool
}

func DefaultOptions() *ClientOptions {
	return &ClientOptions{
		AutoReconnect: true,
		AutoLogin:     true,
	}
}

//Client represents a connection to a deepstream.io server
type Client struct {
	AuthParams map[string]interface{}
	Connector  *Connector
	Options    *ClientOptions
}

//New creates a new client
func New(url, username, password string, optionsOrNil ...*ClientOptions) (*Client, error) {
	options := DefaultOptions()
	if len(optionsOrNil) == 1 {
		options = optionsOrNil[0]
	}
	authParams := map[string]interface{}{}
	if username != "" {
		authParams = map[string]interface{}{
			"username": username,
			"password": password,
		}
	}
	cli := &Client{
		Connector:  NewConnector(url),
		Options:    options,
		AuthParams: authParams,
	}

	cli.Connector.AddMessageHandler(cli.OnMessage)

	err := cli.StartMonitoringConnection()
	if err != nil {
		return cli, err
	}

	return cli, nil
}

func (c *Client) StartMonitoringConnection() error {
	err := c.Connector.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) OnMessage(msg *message.Message) {
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
	if c.Options.AutoLogin {
		return c.Login()
	}
	return nil
}

func (c *Client) Login() error {
	state := c.GetConnectionState()
	if state != interfaces.ConnectionStateAwaitingConnection {
		return fmt.Errorf("The connection should be restored before logging in (%s).", state)
	}

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
	default:
		fmt.Println("Message not understood!")
	}

	return nil
}

func (c *Client) handleAuthenticationAck(msg *message.Message) error {
	c.Connector.ConnectionState = interfaces.ConnectionStateOpen
	return nil
}

//Connect with deepstream.io server
//func (c *Client) Connect() error {
//c.ConnectionState = interfaces.ConnectionStateAwaitingConnection

//err := c.Protocol.Connect()
//if err != nil {
//return c.Error(err)
//}

////Send Challenge Response
//err = c.sendChallengeResponse()
//if err != nil {
//return c.Error(err)
//}

//err = c.receiveAck("C")
//if err != nil {
//return c.Error(err)
//}

//return nil
//}

//func (c *Client) sendChallengeResponse() error {
//challenge := message.NewChallengeResponseAction(c.URL)
//return c.Protocol.SendAction(challenge)
//}

//func (c *Client) receiveAck(expectedTopic string) error {
//// Receive connection Ack
//actions, err := c.Protocol.RecvActions()
//if err != nil {
//return err
//}

//if len(actions) != 1 {
////TODO: change this
//return fmt.Errorf("Expected Ack")
//}

//action := actions[0]
//if a, ok := action.(*message.AckAction); !ok || a.Topic != expectedTopic {
//return fmt.Errorf("Expected Ack with topic %s.", expectedTopic)
//}

//return nil
//}

////Close connection to deepstream.io server
//func (c *Client) Close() error {
//err := c.Protocol.Close()
//if err != nil {
//return c.Error(err)
//}

//c.ConnectionState = interfaces.ConnectionStateClosed
//return nil
//}

////Login with deepstream.io server
//func (c *Client) Login(authParams map[string]interface{}) error {
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

//return nil
//}

////Error handlers errors in client
//func (c *Client) Error(err error) error {
//c.ConnectionState = interfaces.ConnectionStateError
//return err
//}
