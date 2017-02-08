package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/heynemann/deepstream.io-client-go/deepstream"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
)

var receivedEvents []*deepstream.EventMessage
var receivedErrors []error

func handleClientErrors(err error) {
	receivedErrors = append(receivedErrors, err)
}

func theClientIsInitialised() error {
	var err error
	if _, ok := deepstream.TestServers[defaultPort]; !ok {
		err = deepstream.StartTestServer(defaultPort)
		if err != nil {
			return err
		}
	}
	opts := deepstream.DefaultOptions()
	opts.AutoLogin = false
	opts.HeartbeatIntervalMs = 100000000
	opts.ErrorHandler = handleClientErrors
	url := fmt.Sprintf("localhost:%d", defaultPort)
	client, err = deepstream.New(url, opts)
	if err != nil {
		return err
	}
	return nil
}

func theClientIsInitialisedWithASmallHeartbeatInterval() error {
	var err error
	opts := deepstream.DefaultOptions()
	opts.AutoLogin = false
	opts.HeartbeatIntervalMs = 300
	opts.ErrorHandler = handleClientErrors
	url := fmt.Sprintf("localhost:%d", defaultPort)
	client, err = deepstream.New(url, opts)
	if err != nil {
		return err
	}
	return nil
}

func theClientsConnectionStateIs(arg1 string) error {
	state := client.GetConnectionState()
	if string(state) != arg1 {
		return fmt.Errorf("Expected state to be %s but it was %s.", arg1, state)
	}
	return nil
}

func reconnectServerAndCheckState() error {
	err := deepstream.StartTestServer(defaultPort)
	if err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	state := client.GetConnectionState()
	if string(state) != string(interfaces.ConnectionStateReconnecting) {
		return fmt.Errorf("Expected state to be RECONNECTING but it was %s.", state)
	}
	return nil
}

func theClientLogsInWithUsernameAndPassword(username, password string) error {
	client.Options.Username = username
	client.Options.Password = password
	client.AuthParams = map[string]interface{}{
		"username": username,
		"password": password,
	}

	client.Login()

	time.Sleep(20 * time.Millisecond)
	return nil
}

func theClientThrowsAnErrorWithMessage(expectedError, expectedMessage string) error {
	for _, err := range receivedErrors {
		if err.Error() == expectedMessage {
			return nil
		}
	}

	errors := make([]string, len(receivedErrors))
	for i, err := range receivedErrors {
		errors[i] = err.Error()
	}

	return fmt.Errorf(
		"The error with message '%s' did not happen (errors: %s).",
		expectedMessage,
		strings.Join(errors, ", "),
	)
}

func theClientThrowsAnErrorMessage(expectedMessage string) error {
	for _, err := range receivedErrors {
		if err.Error() == expectedMessage {
			return nil
		}
	}
	return fmt.Errorf("The error with message '%s' did not happen.", expectedMessage)
}

func theClientIsOnTheSecondServer() error {
	if client.Connector.URL != "localhost:9998" {
		return fmt.Errorf("Client should be connected to second server but it is connected to %s", client.Connector.URL)
	}
	return nil
}

func theLastLoginWasSuccessful() error {
	if len(receivedErrors) > 0 {
		errors := make([]string, len(receivedErrors))
		for i, err := range receivedErrors {
			errors[i] = err.Error()
		}
		return fmt.Errorf(
			"The login was not successful. The current connection state is %s (errors: %s).",
			client.GetConnectionState(),
			strings.Join(errors, ", "),
		)
	}
	return nil
}

func handleEventReceived(msg *deepstream.EventMessage) error {
	receivedEvents = append(receivedEvents, msg)
	return nil
}

func theClientSubscribesToAnEventNamed(eventName string) error {
	err := client.Event.Subscribe(eventName, handleEventReceived)
	time.Sleep(10 * time.Millisecond)
	return err
}

func theClientUnsubscribesFromAnEventNamed(eventName string) error {
	err := client.Event.Unsubscribe(eventName)
	time.Sleep(10 * time.Millisecond)
	return err
}

func theClientReceivedEvent(eventName, eventData string) error {
	for _, msg := range receivedEvents {
		if msg.Event == eventName && msg.Data[0] == eventData {
			return nil
		}
	}

	return fmt.Errorf(
		"None of the received events (%s) match the event %s with data %s.",
		receivedEvents,
		eventName,
		eventData,
	)
}

func theClientPublishesAnEvent(eventName, eventData string) error {
	err := client.Event.Publish(eventName, eventData)
	time.Sleep(10 * time.Millisecond)
	return err
}

func theClientListensToEvents(eventPattern string) error {
	err := client.Event.Listen(eventPattern, handleEventReceived)
	time.Sleep(10 * time.Millisecond)
	return err
}
