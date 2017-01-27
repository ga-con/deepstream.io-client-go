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
	Connection      interfaces.Connection
	ConnectionState interfaces.ConnectionState
}

func loadConnection(url string, options []interface{}) interfaces.Connection {
	var conn interfaces.Connection

	if len(options) >= 1 {
		conn = options[1].(interfaces.Connection)
	}

	if conn == nil {
		var err error
		conn, err = NewWebsocketConnection(url)
		if err != nil {
			return err
		}
	}

	return conn
}

//New creates a new client
func New(url string, options ...interface{}) (*Client, error) {
	connection := loadConnection(options)

	cli := &Client{
		URL:             url,
		ConnectionState: interfaces.ConnectionStateClosed,
		Connection:      connection,
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

	err := c.Connection.Connect()
	if err != nil {
		return c.Error(err)
	}

	//Send Challenge Response
	//err = c.sendChallengeResponse()
	//if err != nil {
	//return c.Error(err)
	//}

	//err = c.receiveAck("C")
	//if err != nil {
	//return c.Error(err)
	//}

	return nil
}

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

//Close connection to deepstream.io server
func (c *Client) Close() error {
	//err := c.Protocol.Close()
	//if err != nil {
	//return c.Error(err)
	//}

	c.ConnectionState = interfaces.ConnectionStateClosed
	return nil
}

//Login with deepstream.io server
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

//Error handlers errors in client
func (c *Client) Error(err error) error {
	c.ConnectionState = interfaces.ConnectionStateError
	return err
}
