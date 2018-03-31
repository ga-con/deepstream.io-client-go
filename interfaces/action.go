// deepstream.io-client-go
// https://github.com/ga-con/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package interfaces

//Action represents a single action coming from deepstream.io
type Action interface {
	//ToAction returns the action in raw format that deepstream.io understands
	ToAction() string
}
