// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package testing

//MockProtocol should be used for unit tests
type MockProtocol struct {
	Error           error
	IsAuthenticated bool
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

//Authenticate mock protocol
func (m *MockProtocol) Authenticate() error {
	if m.Error != nil {
		return m.Error
	}

	m.IsAuthenticated = true

	return nil
}
