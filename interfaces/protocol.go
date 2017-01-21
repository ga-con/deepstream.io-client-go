// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package interfaces

//Protocol specifies the transport protocol for the client
type Protocol interface {
	Connect() error
	Authenticate(authParams map[string]interface{}) error
}
