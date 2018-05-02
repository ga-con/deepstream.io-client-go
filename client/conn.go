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
	"time"
	"sync"
	"github.com/gorilla/websocket"
	"github.com/ga-con/deepstream.io-client-go/interfaces"
	"github.com/ga-con/deepstream.io-client-go/message"
	"github.com/jpillora/backoff"
	"math/rand"
	"log"
)

type AuthUser struct {
	Username string
	Password string
	Token    string
}

// ClientOption is a function on the options for a connection.
type ClientOption func(*ClientOptions) error

type ClientOptions struct {
	// RecIntvlMin specifies the initial reconnecting interval,
	// default to 2 seconds
	RecIntvlMin time.Duration
	// RecIntvlMax specifies the maximum reconnecting interval,
	// default to 30 seconds
	RecIntvlMax time.Duration
	// RecIntvlFactor specifies the rate of increase of the reconnection
	// interval, default to 1.5
	RecIntvlFactor float64
	// HandshakeTimeout specifies the duration for the handshake to complete,
	// default to 2 seconds
	HandshakeTimeout time.Duration

	AuthUser AuthUser
}

// GetDefaultOptions returns default configuration options for the client.
func GetDefaultOptions() ClientOptions {
	return ClientOptions{
		RecIntvlMin:      2 * time.Second,
		RecIntvlMax:      30 * time.Second,
		RecIntvlFactor:   1.5,
		HandshakeTimeout: 2 * time.Second,
	}
}

//Client represents a connection to a deepstream.io server
type Client struct {
	Options         ClientOptions
	URL             string
	ConnectionState interfaces.ConnectionState
	mu              sync.Mutex
	dialErr         error
	isConnected     bool
	IsLogin         bool
	dialer          *websocket.Dialer
	*websocket.Conn
}

//Dial creates a new client connection.
func Dial(url string, options ...ClientOption) (*Client, error) {
	opts := GetDefaultOptions()
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}

	cli := &Client{
		URL:             fmt.Sprintf("ws://%s/deepstream", url),
		ConnectionState: interfaces.ConnectionStateClosed,
		dialer:          websocket.DefaultDialer,
		Options:         opts,
	}
	cli.dialer.HandshakeTimeout = cli.Options.HandshakeTimeout

	go func() {
		cli.connect()
	}()

	// wait on first attempt
	time.Sleep(cli.Options.HandshakeTimeout)

	return cli, nil
}

func (cli *Client) connect() {
	b := &backoff.Backoff{
		Min:    cli.Options.RecIntvlMin,
		Max:    cli.Options.RecIntvlMax,
		Factor: cli.Options.RecIntvlFactor,
		Jitter: true,
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for {
		nextItvl := b.Duration()

		wsConn, _, err := cli.dialer.Dial(cli.URL, nil)

		cli.mu.Lock()
		cli.Conn = wsConn
		cli.dialErr = err
		cli.isConnected = err == nil
		cli.mu.Unlock()

		if err == nil {
			log.Printf("Dial: connection was successfully established with %s\n", cli.URL)

			err = cli.getAuthChallenge()
			if err != nil {
				fmt.Println(err)
				return
			}

			//Send Challenge Response
			err = cli.sendChallengeResponse()
			if err != nil {
				fmt.Println(err)

				return
			}

			err = cli.receiveAck("C")
			if err != nil {
				fmt.Println(err)

				return
			}

			var param map[string]interface{}
			if len(cli.Options.AuthUser.Token) > 0 {
				param = map[string]interface{}{"token": cli.Options.AuthUser.Token}
			} else {
				param = map[string]interface{}{
					"username": cli.Options.AuthUser.Username,
					"password": cli.Options.AuthUser.Password}
			}

			if err := cli.Login(param); err != nil {
				log.Println(err)
			} else {
				cli.IsLogin = true
			}

			break
		} else {
			log.Println("Dial: will try again in", nextItvl, "seconds.", " Error:", err)
		}

		time.Sleep(nextItvl)
	}
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
// Close closes the underlying network connection without
// sending or waiting for a close frame.
func (cli *Client) Close() error {
	cli.mu.Lock()
	defer cli.mu.Unlock()

	if cli.Conn != nil {
		cli.Conn.Close()
	}
	cli.isConnected = false
	cli.IsLogin = false

	cli.ConnectionState = interfaces.ConnectionStateClosed
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

			if err != nil {
				fmt.Errorf("handlerConnection:%v", err)
				return
			}

			for _, act := range acts {
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
	if c.IsConnected() {
		c.mu.Lock()
		defer c.mu.Unlock()

		msg := action.ToAction()
		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			c.closeAndRecconect()
			return err
		}
		return nil
	}
	return fmt.Errorf("Not Connected!")
}

//RecvActions receives actions from the websocket stream
func (c *Client) RecvActions() ([]interfaces.Action, error) {
	if c.IsConnected() {
		_, body, err := c.Conn.ReadMessage()
		if err != nil {
			c.closeAndRecconect()
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

	return []interfaces.Action{}, fmt.Errorf("Not Connected!")

}

// IsConnected returns the WebSocket connection state
func (cli *Client) IsConnected() bool {
	cli.mu.Lock()
	defer cli.mu.Unlock()

	return cli.isConnected
}

// CloseAndRecconect will try to reconnect.
func (rc *Client) closeAndRecconect() {
	rc.Close()
	go func() {
		rc.connect()
	}()

}
