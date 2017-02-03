// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream

import (
	"encoding/json"
	"fmt"

	"github.com/heynemann/deepstream.io-client-go/interfaces"
)

//ChallengeResponseAction creates the action to respond to the server's challenge with
type ChallengeResponseAction struct {
	URL string
}

//NewChallengeResponseAction creates a new message
func NewChallengeResponseAction(url string) *ChallengeResponseAction {
	return &ChallengeResponseAction{
		URL: url,
	}
}

//ToAction converts to the action to be sent to the server
func (a *ChallengeResponseAction) ToAction() string {
	return fmt.Sprintf(
		"C%sCHR%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.URL,
		interfaces.MessageSeparator,
	)
}

//AuthRequestAction represents an action to submit an authentication request to the server
type AuthRequestAction struct {
	AuthParams string
}

//NewAuthRequestAction creates a new auth request action
func NewAuthRequestAction(authParams map[string]interface{}) (*AuthRequestAction, error) {
	data, err := json.Marshal(authParams)
	if err != nil {
		return nil, err
	}
	return &AuthRequestAction{AuthParams: string(data)}, nil
}

//ToAction converts to the action to be sent to the server
func (a *AuthRequestAction) ToAction() string {
	return fmt.Sprintf(
		"A%sREQ%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.AuthParams,
		interfaces.MessageSeparator,
	)
}

type SubscribeToEventAction struct {
	Event string
}

func (a *SubscribeToEventAction) ToAction() string {
	return fmt.Sprintf(
		"E%sS%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.Event,
		interfaces.MessageSeparator,
	)
}
