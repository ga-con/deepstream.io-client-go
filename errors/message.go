// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package errors

import "errors"

var (
	//ErrEmptyRawMessage error
	ErrEmptyRawMessage = errors.New("Message can't be parsed since it's empty and does not conform to the deepstream.io spec")
)
