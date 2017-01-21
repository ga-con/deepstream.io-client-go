// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package message

import (
	"fmt"

	"github.com/heynemann/deepstream.io-client-go/interfaces"
)

//ChallengeAction represents a challenge action coming from the server
type ChallengeAction struct {
	Message
}

//NewChallengeAction creates a challenge action from a raw string
func NewChallengeAction(msg *Message) (*ChallengeAction, error) {
	action := &ChallengeAction{*msg}

	return action, nil
}

func (a *ChallengeAction) ToAction() string {
	return fmt.Sprintf(
		"C%sCH%s",
		interfaces.MessagePartSeparator,
		interfaces.MessageSeparator,
	)
}

type ChallengeResponseAction struct {
	URL string
}

func NewChallengeResponseAction(url string) *ChallengeResponseAction {
	return &ChallengeResponseAction{
		URL: url,
	}
}

func (a *ChallengeResponseAction) ToAction() string {
	return fmt.Sprintf(
		"C%sCHR%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.URL,
		interfaces.MessageSeparator,
	)
}

type AckAction struct {
	Message
}

func NewAckAction(msg *Message) (*AckAction, error) {
	return &AckAction{*msg}, nil
}

func (a *AckAction) ToAction() string {
	return fmt.Sprintf(
		"C%sA%s",
		interfaces.MessagePartSeparator,
		interfaces.MessageSeparator,
	)
}

type AuthRequestAction struct {
	AuthParams string
}

func NewAuthRequestAction(msg *Message) (*AuthRequestAction, error) {
	authParams := ""
	if len(msg.RawData) > 0 {
		authParams = msg.RawData[0]
	}
	return &AuthRequestAction{AuthParams: authParams}, nil
}

func (a *AuthRequestAction) ToAction() string {
	return fmt.Sprintf(
		"A%sREQ%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.AuthParams,
		interfaces.MessageSeparator,
	)
}
