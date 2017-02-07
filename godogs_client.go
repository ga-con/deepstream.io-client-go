package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/heynemann/deepstream.io-client-go/deepstream"
)

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
	url := fmt.Sprintf("127.0.0.1:%d", defaultPort)
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
	url := fmt.Sprintf("127.0.0.1:%d", defaultPort)
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
	return fmt.Errorf("The error with message '%s' did not happen.", expectedMessage)
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
	if client.Connector.URL != "127.0.0.1:9998" {
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

func handleEventReceived() {

}

func theClientSubscribesToAnEventNamed(eventName string) error {
	return client.Event.Subscribe(eventName, handleEventReceived)
}
