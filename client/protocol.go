// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
	"github.com/heynemann/deepstream.io-client-go/message"
)

//DeepstreamProtocol is the default protocol for deepstream.io
type DeepstreamProtocol struct {
	Connection interfaces.Connection
}

//NewDeepstreamProtocol creates a new instance
func NewDeepstreamProtocol(connection interfaces.Connection) (*DeepstreamProtocol, error) {
	dsp := &DeepstreamProtocol{
		Connection: connection,
	}
	return dsp, nil
}

//Connect to deepstream.io
//func (w *DeepstreamProtocol) Connect() error {
//url := fmt.Sprintf("ws://%s/deepstream", w.URL)
//c, _, err := websocket.DefaultDialer.Dial(url, nil)
//if err != nil {
//return err
//}

//w.Client = c

//err = w.getAuthChallenge()
//if err != nil {
//return err
//}
//return nil
//}

//GetAuthChallenge receives the message from deepstream ensuring the auth challenge has been met
func (w *DeepstreamProtocol) GetAuthChallenge() error {
	body, err := w.Connection.ReceiveMessage()
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
func (w *DeepstreamProtocol) Close() error {
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
func (w *DeepstreamProtocol) SendAction(action interfaces.Action) error {
	msg := action.ToAction()
	err := w.Client.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

//RecvActions receives actions from the websocket stream
func (w *DeepstreamProtocol) RecvActions() ([]interfaces.Action, error) {
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
