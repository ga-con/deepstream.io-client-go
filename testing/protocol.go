// deepstream.io-client-go
// https://github.com/ga-con/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package testing

import "github.com/ga-con/deepstream.io-client-go/interfaces"

//MockProtocol should be used for unit tests
type MockProtocol struct {
	Error           error
	IsClosed        bool
	HasConnected    bool
	IsAuthenticated bool
	AuthParams      map[string]interface{}
}

//NewMockProtocol returns a new MockProtocol
func NewMockProtocol(errOrNil ...error) *MockProtocol {
	var err error
	if len(errOrNil) == 1 {
		err = errOrNil[0]
	}
	return &MockProtocol{
		Error: err,
	}
}

//Connect mocks connection
func (m *MockProtocol) Connect() error {
	if m.Error != nil {
		return m.Error
	}

	m.HasConnected = true
	return nil
}

func (m *MockProtocol) SendAction(action interfaces.Action) error {
	return nil
}

func (m *MockProtocol) RecvActions() ([]interfaces.Action, error) {
	return nil, nil
}

//Authenticate mock protocol
func (m *MockProtocol) Authenticate(authParams map[string]interface{}) error {
	m.AuthParams = authParams

	if m.Error != nil {
		return m.Error
	}

	m.IsAuthenticated = true

	return nil
}

//Close mock connection
func (m *MockProtocol) Close() error {
	if m.Error != nil {
		return m.Error
	}

	m.IsClosed = true
	return nil
}
