// deepstream.io-client-go
// https://github.com/ga-con/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package errors

import "errors"

var (
	//ErrUnknownAction error
	ErrUnknownAction = errors.New("Action with the specified type could not be understood.")
)
