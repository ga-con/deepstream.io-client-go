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
	"strconv"
	"strings"

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
	data := fmt.Sprintf(`{"username":"%s","password":"%s"}`, authParams["username"], authParams["password"])
	return &AuthRequestAction{AuthParams: data}, nil
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

//SubscribeToEventAction sends a message to deepstream to subscribe to a given event
type SubscribeToEventAction struct {
	Event string
}

//ToAction converts to the action to be sent to the server
func (a *SubscribeToEventAction) ToAction() string {
	return fmt.Sprintf(
		"E%sS%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.Event,
		interfaces.MessageSeparator,
	)
}

//UnsubscribeFromEventAction sends a message to deepstream to unsubscribe to a given event
type UnsubscribeFromEventAction struct {
	Event string
}

//ToAction converts to the action to be sent to the server
func (a *UnsubscribeFromEventAction) ToAction() string {
	return fmt.Sprintf(
		"E%sUS%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.Event,
		interfaces.MessageSeparator,
	)
}

//PublishEventAction sends a message to deepstream to publish an event
type PublishEventAction struct {
	Event string
	Data  []interface{}
}

func serializeData(source []interface{}) ([]string, error) {
	data := make([]string, len(source))

	for i, item := range source {
		if item == nil {
			data[i] = "L"
			continue
		}

		switch t := item.(type) {
		default:
			marsh, err := json.Marshal(t)
			if err != nil {
				return nil, err
			}
			data[i] = fmt.Sprintf("O%s", marsh)
		case string:
			data[i] = fmt.Sprintf("S%s", t)
		case bool:
			data[i] = "F"
			if t {
				data[i] = "T"
			}
		case int, int16, int32, int64, uint, float32, float64:
			data[i] = fmt.Sprintf("N%d", t)
		case *string:
			data[i] = fmt.Sprintf("S%s", *t)
		case *bool:
			data[i] = "F"
			if *t {
				data[i] = "T"
			}
		}
	}

	return data, nil
}

func deserializeData(source []string) ([]interface{}, error) {
	data := make([]interface{}, len(source))

	for i, item := range source {
		switch item[0] {
		case 'L':
			data[i] = nil
		case 'O':
			var obj interface{}
			err := json.Unmarshal([]byte(item[1:]), &obj)
			if err != nil {
				return nil, err
			}
			data[i] = obj
		case 'S':
			data[i] = item[1:]
		case 'F':
			data[i] = false
		case 'T':
			data[i] = true
		case 'N':
			num, err := strconv.ParseFloat(item[1:], 64)
			if err != nil {
				return nil, err
			}
			data[i] = num
		}
	}

	return data, nil
}

//ToAction converts to the action to be sent to the server
func (a *PublishEventAction) ToAction() (string, error) {
	data, err := serializeData(a.Data)
	if err != nil {
		return "", err
	}

	dataStr := strings.Join(data, interfaces.MessagePartSeparator)

	return fmt.Sprintf(
		"E%sEVT%s%s%s%s%s",
		interfaces.MessagePartSeparator,
		interfaces.MessagePartSeparator,
		a.Event,
		interfaces.MessagePartSeparator,
		dataStr,
		interfaces.MessageSeparator,
	), nil
}

//PongAction sends a response to a ping request
type PongAction struct{}

//ToAction converts to the action to be sent to the server
func (a *PongAction) ToAction() string {
	return fmt.Sprintf(
		"C%sPO%s",
		interfaces.MessagePartSeparator,
		interfaces.MessageSeparator,
	)
}
