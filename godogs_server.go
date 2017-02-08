package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/heynemann/deepstream.io-client-go/deepstream"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
)

func theTestServerIsReady() error {
	var err error
	err = deepstream.StartTestServer(defaultPort)
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	return nil
}

func theServerHasActiveConnections(expectedConnections int) error {
	activeConnections := deepstream.TestServers[defaultPort].ActiveConnections
	if activeConnections != expectedConnections {
		return fmt.Errorf("Expected %d active connections to server, but %d are active.", expectedConnections, activeConnections)
	}
	return nil
}

func theServerSendsTheMessage(topic, action string, data ...string) func() error {
	return func() error {
		dataMsg := ""
		if len(data) > 0 {
			dataMsg = fmt.Sprintf(
				"%s%s", interfaces.MessagePartSeparator,
				strings.Join(data, interfaces.MessagePartSeparator),
			)
		}

		var err error
		for _, ts := range deepstream.TestServers {
			err = ts.SendMessage(
				fmt.Sprintf(
					"%s%s%s%s%s", topic, interfaces.MessagePartSeparator,
					action, dataMsg, interfaces.MessageSeparator,
				),
			)
		}

		time.Sleep(20 * time.Millisecond)

		return err
	}
}

func theServerReceivedTheMessage(topic, action string, data ...string) func() error {
	return func() error {
		dataMsg := ""

		if len(data) > 0 {
			dataMsg = fmt.Sprintf(
				"%s%s",
				interfaces.MessagePartSeparator,
				strings.Join(data, interfaces.MessagePartSeparator),
			)
		}

		var err error
		expectedMessage := fmt.Sprintf(
			"%s%s%s%s%s", topic, interfaces.MessagePartSeparator,
			action, dataMsg, interfaces.MessageSeparator,
		)
		for _, ts := range deepstream.TestServers {
			err = ts.HasMessage(expectedMessage)
			if err == nil {
				return nil
			}
		}
		return err
	}
}

func theSecondTestServerIsReady() error {
	var err error
	err = deepstream.StartTestServer(9998)
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	return nil
}

func theSecondServerHasActiveConnections(expectedConnections int) error {
	activeConnections := deepstream.TestServers[9998].ActiveConnections
	if activeConnections != expectedConnections {
		return fmt.Errorf("Expected %d active connections to server, but %d are active.", expectedConnections, activeConnections)
	}
	return nil
}

func theServerHasReceivedMessages(numberOfMessages int) error {
	recMessages := len(deepstream.TestServers[defaultPort].ReceivedMessages)
	if recMessages != numberOfMessages {
		return fmt.Errorf(
			"Expected server to have received %d messages, but there were %d messages.",
			numberOfMessages,
			recMessages,
		)
	}
	return nil
}

func theConnectionToTheServerIsLost() error {
	deepstream.TestServers[defaultPort].Stop()
	delete(deepstream.TestServers, defaultPort)
	time.Sleep(50 * time.Millisecond)
	return nil
}
