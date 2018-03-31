// deepstream.io-client-go
// https://github.com/ga-con/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/ga-con/deepstream.io-client-go/interfaces"
	"github.com/ga-con/deepstream.io-client-go/message"
)

//WebsocketProtocol is the default protocol for deepstream.io
type WebsocketProtocol struct {
	URL    string
	Client *websocket.Conn
}

//NewWebsocketProtocol creates a new instance
func NewWebsocketProtocol(url string) (*WebsocketProtocol, error) {
	ws := &WebsocketProtocol{
		URL: url,
	}
	return ws, nil
}

//Connect to deepstream.io
func (w *WebsocketProtocol) Connect() error {
	url := fmt.Sprintf("ws://%s/deepstream", w.URL)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	w.Client = c

	err = w.getAuthChallenge()
	if err != nil {
		return err
	}
	return nil
}

func (w *WebsocketProtocol) getAuthChallenge() error {
	_, body, err := w.Client.ReadMessage()
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

//SendAction writes an action in the websocket stream
func (w *WebsocketProtocol) SendAction(action interfaces.Action) error {
	msg := action.ToAction()
	err := w.Client.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

//RecvActions receives actions from the websocket stream
func (w *WebsocketProtocol) RecvActions() ([]interfaces.Action, error) {
	_, body, err := w.Client.ReadMessage()
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
