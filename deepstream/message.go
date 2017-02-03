// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream

import (
	"strings"

	"github.com/heynemann/deepstream.io-client-go/errors"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
)

//Message represents a message received from deepstream.io
type Message struct {
	Raw    string
	Topic  string
	Action string
	Data   []string
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
	m.Data = parts[2:]

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
