// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package testing

import "fmt"

//MockConnection should be used for unit tests
type MockConnection struct {
	Error                 error
	IsClosed              bool
	HasConnected          bool
	MessagesSent          [][]byte
	MessagesBeingReceived [][]byte
	MessageIndex          int
}

//NewMockConnection returns a new MockConnection
func NewMockConnection(errOrNil ...error) *MockConnection {
	var err error
	if len(errOrNil) == 1 {
		err = errOrNil[0]
	}
	return &MockConnection{
		Error:                 err,
		MessagesSent:          [][]byte{},
		MessagesBeingReceived: [][]byte{},
	}
}

//Connect mocks connection
func (m *MockConnection) Connect() error {
	if m.Error != nil {
		return m.Error
	}

	m.HasConnected = true
	return nil
}

//SendMessage mocks sending a message
func (m *MockConnection) SendMessage(msg []byte) error {
	if m.Error != nil {
		return m.Error
	}

	m.MessagesSent = append(m.MessagesSent, msg)
	return nil
}

//ReceiveMessage mocks receiving a message
func (m *MockConnection) ReceiveMessage() ([]byte, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	if m.MessageIndex > len(m.MessagesBeingReceived) {
		return nil, fmt.Errorf("No message could be found to receive. Please use MockConnection.SetNextReceiveMessage.")
	}

	msg := m.MessagesBeingReceived[m.MessageIndex]
	m.MessageIndex++

	return msg, nil
}

//SetNextReceiveMessage to specify how the mock will behave
func (m *MockConnection) SetNextReceiveMessage(msg []byte) {
	m.MessagesBeingReceived = append(m.MessagesBeingReceived, msg)
}

//Close mock connection
func (m *MockConnection) Close() error {
	if m.Error != nil {
		return m.Error
	}

	m.IsClosed = true
	return nil
}
