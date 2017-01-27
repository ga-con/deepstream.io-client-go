// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package interfaces

//BinaryMessage represents a binary message
type BinaryMessage struct {
	MessageType int
	Payload     []byte
	Err         error
}

//Connection to the deepstream server
type Connection interface {
	Connect() error
	SendMessage(msg []byte) error
	ReceiveMessage() ([]byte, error)
}
