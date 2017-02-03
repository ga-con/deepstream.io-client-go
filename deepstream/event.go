// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream

import "fmt"

//EventHandler is the function type for all event handlers
type EventHandler func()

//EventSubscription represents a subscription to a given event
type EventSubscription struct {
	Event    string
	Handlers []EventHandler
	Acked    bool
}

//EventManager handles all event related operations
type EventManager struct {
	client        *Client
	Subscriptions map[string]*EventSubscription
}

//NewEventManager creates a new Event Manager
func NewEventManager(cli *Client) *EventManager {
	return &EventManager{
		client:        cli,
		Subscriptions: map[string]*EventSubscription{},
	}
}

//Subscribe to events in deepstream.io
func (e *EventManager) Subscribe(event string, handler EventHandler) error {
	if sub, ok := e.Subscriptions[event]; ok {
		sub.Handlers = append(sub.Handlers, handler)
		return nil
	}

	e.Subscriptions[event] = &EventSubscription{
		Acked:    false,
		Event:    event,
		Handlers: []EventHandler{handler},
	}

	action := &SubscribeToEventAction{
		Event: event,
	}
	return e.client.Connector.WriteMessage([]byte(action.ToAction()))
}

func (e *EventManager) handleEventSubscriptionAck(msg *Message) error {
	if len(msg.Data)%2 != 0 {
		return fmt.Errorf("Invalid data returned for event acknowledge: %v", msg.Data)
	}

	for i := 0; i < len(msg.Data); i += 2 {
		flag := msg.Data[0]
		event := msg.Data[1]

		if flag != "S" {
			return fmt.Errorf("Invalid subscription acknowledge for event subscription: %s", flag)
		}

		if sub, ok := e.Subscriptions[event]; ok {
			sub.Acked = true
		} else {
			return fmt.Errorf("Received subscription confirmation for unknown subscription: %s", event)
		}
	}
	return nil
}
