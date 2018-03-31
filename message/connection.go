// deepstream.io-client-go
// https://github.com/ga-con/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package message

import (
	"fmt"
	"github.com/ga-con/deepstream.io-client-go/interfaces"
	"strings"
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

type CreateOrReadAction struct {
	Message
}

func NewCreateOrReadAction(msg *Message) (*CreateOrReadAction, error) {
	return &CreateOrReadAction{*msg}, nil
}

func (a *CreateOrReadAction) ToAction() string {
	return fmt.Sprintf(
		"R%sCR%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.RawData[0],
		interfaces.MessageSeparator,
	)
}

type UpdateAction struct {
	Message
}

func NewUpdateAction(msg *Message) (*UpdateAction, error) {
	return &UpdateAction{*msg}, nil
}

// 	When the server sends the message R|U|subscribeRecord|125|{"name":"Smith","pets":[{"name":"Ruffus","type":"dog","age":1}]}+
func (a *UpdateAction) ToAction() string {
	return fmt.Sprintf(
		"R%sU%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		strings.Join(a.RawData, interfaces.MessagePartSeparator),
		interfaces.MessageSeparator,
	)
}

type PathAction struct {
	Message
}

func NewPathAction(msg *Message) (*PathAction, error) {
	return &PathAction{*msg}, nil
}

func (a *PathAction) ToAction() string {
	return fmt.Sprintf(
		"R%sP%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		strings.Join(a.RawData, interfaces.MessagePartSeparator),
		interfaces.MessageSeparator,
	)
}

type ReadAction struct {
	Message
}

func NewReadAction(msg *Message) (*ReadAction, error) {
	return &ReadAction{*msg}, nil
}

func (a *ReadAction) ToAction() string {
	return fmt.Sprintf(
		"R%sR%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		strings.Join(a.RawData, interfaces.MessagePartSeparator),
		interfaces.MessageSeparator,
	)
}

type EventAction struct {
	Message
}

func NewEventAction(msg *Message) (*EventAction, error) {
	return &EventAction{*msg}, nil
}

//  E|EVT|test1|SyetAnotherValue+
func (a *EventAction) ToAction() string {
	return fmt.Sprintf(
		"E%sEVT%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		strings.Join(a.RawData, interfaces.MessagePartSeparator),
		interfaces.MessageSeparator,
	)
}

// E|S|test1+
type SubscribeAction struct {
	Message
}

func NewSubscribeAction(msg *Message) (*SubscribeAction, error) {
	return &SubscribeAction{*msg}, nil
}

func (a *SubscribeAction) ToAction() string {
	return fmt.Sprintf(
		"E%sS%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		strings.Join(a.RawData, interfaces.MessagePartSeparator),
		interfaces.MessageSeparator,
	)
}

type PingAction struct {
	Message
}

func NewPingAction(msg *Message) (*PingAction, error) {
	return &PingAction{*msg}, nil
}

func (a *PingAction) ToAction() string {
	return fmt.Sprintf(
		"C%sPI",
		interfaces.MessagePartSeparator,
	)
}

type PongAction struct {
	Message
}

func NewPongAction(msg *Message) (*PongAction, error) {
	return &PongAction{*msg}, nil
}

func (a *PongAction) ToAction() string {
	return fmt.Sprintf(
		"C%sPO",
		interfaces.MessagePartSeparator,
	)
}
