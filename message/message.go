// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package message

import (
	"strings"

	"github.com/heynemann/deepstream.io-client-go/errors"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
)

var (
	//AvailableMessageTypes returns all the available message types
	AvailableMessageTypes = map[string]func(*Message) (interfaces.Action, error){
		interfaces.ActionChallenge: func(msg *Message) (interfaces.Action, error) { return NewChallengeAction(msg) },
		interfaces.ActionAck:       func(msg *Message) (interfaces.Action, error) { return NewAckAction(msg) },
	}
)

//Data represents a portion of data coming from client
type Data struct {
	Type  interfaces.DataType
	Value interface{}
}

//Message represents a message received from deepstream.io
type Message struct {
	Raw     string
	Topic   string
	Action  string
	RawData []string
	Data    []Data
}

//NewMessage creates a new message
func NewMessage(raw string) (*Message, error) {
	msg := &Message{
		Raw: raw,
	}
	err := msg.Parse()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

//Parse the raw message
func (m *Message) Parse() error {
	if m.Raw == "" {
		return errors.ErrEmptyRawMessage
	}

	parts := strings.Split(m.Raw, interfaces.MessagePartSeparator)
	m.Topic = parts[0]
	m.Action = parts[1]
	m.RawData = parts[2:]

	return nil
}

//ParseMessages in a raw string
func ParseMessages(raw string) ([]*Message, error) {
	if raw == "" {
		return nil, errors.ErrEmptyRawMessage
	}

	rawMessages := strings.Split(raw, interfaces.MessageSeparator)
	messages := []*Message{}
	for _, rawMessage := range rawMessages {
		if rawMessage == "" {
			continue
		}
		message, err := NewMessage(rawMessage)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

//CathegorizeAction returns a cathegorized action
func CathegorizeAction(message *Message) (interfaces.Action, error) {
	actionFunc, ok := AvailableMessageTypes[message.Action]
	if !ok {
		return nil, errors.ErrUnknownAction
	}
	action, err := actionFunc(message)
	if err != nil {
		return nil, err
	}
	return action, nil
}
