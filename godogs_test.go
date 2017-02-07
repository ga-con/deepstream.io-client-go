package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/gorilla/websocket"
	"github.com/heynemann/deepstream.io-client-go/deepstream"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
	uuid "github.com/satori/go.uuid"
)

const TimeMultiplier = 0.01

var receivedErrors []error
var client *deepstream.Client
var testServer *TestServer

func afterScenario(interface{}, error) {
	receivedErrors = []error{}
	if client != nil {
		client.Close()
		client = nil
	}

	if testServer != nil {
		testServer.stop()
		time.Sleep(10 * time.Millisecond)
		testServer = nil
	}
}

type TestServer struct {
	Port                 int
	listener             net.Listener
	websocketConnections map[string]*websocket.Conn
	activeConnections    int
	upgrader             websocket.Upgrader
	ReceivedMessages     []string
}

func (ts *TestServer) start() error {
	ts.websocketConnections = map[string]*websocket.Conn{}
	ts.upgrader = websocket.Upgrader{} // use default options
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", ts.Port))
	if err != nil {
		return err
	}
	ts.listener = ln

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/deepstream", ts.handler)

	server := http.Server{
		Handler: serverMux,
	}

	go func() {
		server.Serve(ln)
	}()

	return nil
}

func (ts *TestServer) stop() {
	//close(sendToSocketChan)
	ts.listener.Close()
}

func (ts *TestServer) handler(w http.ResponseWriter, r *http.Request) {
	ts.activeConnections++
	c, err := ts.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade:", err)
		return
	}
	id := uuid.NewV4().String()
	ts.websocketConnections[id] = c

	for {
		_, msg, err := c.ReadMessage()

		if websocket.IsCloseError(err) {
			c.Close()
			delete(ts.websocketConnections, id)
			ts.activeConnections--
			return
		}
		if err != nil {
			delete(ts.websocketConnections, id)
			ts.activeConnections--
			return
		}

		ts.ReceivedMessages = append(ts.ReceivedMessages, string(msg))
	}
}

func (ts *TestServer) sendMessage(message string) error {
	var err error
	for _, conn := range ts.websocketConnections {
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (ts *TestServer) hasMessage(expectedMessage string) error {
	for _, msg := range ts.ReceivedMessages {
		if msg == expectedMessage {
			return nil
		}
	}

	return fmt.Errorf("Message '%s' was not received by the server.", expectedMessage)
}

func startServer(port int) error {
	ts := &TestServer{
		Port: port,
	}
	err := ts.start()
	if err != nil {
		return err
	}

	testServer = ts
	return nil
}

func theTestServerIsReady() error {
	var err error
	err = startServer(9999)
	if err != nil {
		return err
	}
	return nil
}

func theServerHasActiveConnections(expectedConnections int) error {
	if testServer == nil {
		return fmt.Errorf("Server is not up!")
	}

	activeConnections := testServer.activeConnections
	if activeConnections != expectedConnections {
		return fmt.Errorf("Expected %d active connections to server, but %d are active.", expectedConnections, activeConnections)
	}
	return nil
}

func handleClientErrors(err error) {
	receivedErrors = append(receivedErrors, err)
}

func theClientIsInitialised() error {
	var err error
	opts := deepstream.DefaultOptions()
	opts.ErrorHandler = handleClientErrors
	client, err = deepstream.New("127.0.0.1:9999", opts)
	if err != nil {
		return err
	}
	return nil
}

func theClientIsInitialisedWithASmallHeartbeatInterval() error {
	var err error
	opts := deepstream.DefaultOptions()
	opts.AutoLogin = false
	opts.HeartbeatIntervalMs = 20
	opts.ErrorHandler = handleClientErrors
	client, err = deepstream.New("127.0.0.1:9999", opts)
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

func theServerSendsTheMessage(topic, action string) func() error {
	return func() error {
		return testServer.sendMessage(
			fmt.Sprintf("%s%s%s%s", topic, interfaces.MessagePartSeparator, action, interfaces.MessageSeparator),
		)
	}
}

func theClientLogsInWithUsernameAndPassword(username, password string) error {
	client.Options.Username = username
	client.Options.Password = password
	err := client.Login()
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Millisecond)
	return nil
}

func theServerReceivedTheMessage(topic, action string) func() error {
	return func() error {
		expectedMessage := fmt.Sprintf(
			"%s%s%s%s", topic, interfaces.MessagePartSeparator,
			action, interfaces.MessageSeparator,
		)

		return testServer.hasMessage(expectedMessage)
	}
}

func someMsLater(ms int) func() error {
	return func() error {
		duration := int(float64(ms) * TimeMultiplier)
		time.Sleep(time.Duration(duration) * time.Millisecond)
		return nil
	}
}

func theClientThrowsAnErrorWithMessage(expectedError, expectedMessage string) error {
	for _, err := range receivedErrors {
		if err.Error() == expectedMessage {
			return nil
		}
	}
	return fmt.Errorf("The error with message '%s' did not happen.", expectedMessage)
}

func theSecondTestServerIsReady() error {
	return godog.ErrPending
}

func theSecondServerHasActiveConnections(arg1 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageCCH() error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsCCHRFIRSTSERVERURL() error {
	return godog.ErrPending
}

func theServerSendsTheMessageCREJ() error {
	return godog.ErrPending
}

func theServerHasReceivedMessages(arg1 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageCREDSECONDSERVERURL() error {
	return godog.ErrPending
}

func someTimePasses() error {
	return godog.ErrPending
}

func theClientIsOnTheSecondServer() error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsAREQXXXYYY(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastLoginWasSuccessful() error {
	return godog.ErrPending
}

func theServerSendsTheMessageAEINVALIDAUTHDATASinvalidAuthenticationData() error {
	return godog.ErrPending
}

func theLastLoginFailedWithErrorMessage(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageAETOOMANYAUTHATTEMPTSStooManyAuthenticationAttempts() error {
	return godog.ErrPending
}

func theServerResetsItsMessageCount() error {
	return godog.ErrPending
}

func theClientSubscribesToAnEventNamed(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageEAStest(arg1 int) error {
	return godog.ErrPending
}

func theServerReceivedTheMessageEStest(arg1 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsEStest(arg1 int) error {
	return godog.ErrPending
}

func theClientListensToEventsMatching(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsELeventPrefix() error {
	return godog.ErrPending
}

func theServerSendsTheMessageEALeventPrefix() error {
	return godog.ErrPending
}

func theConnectionToTheServerIsLost() error {
	return godog.ErrPending
}

func theClientPublishesAnEventNamedWithData(arg1, arg2 string) error {
	return godog.ErrPending
}

func theServerDidNotRecieveAnyMessages() error {
	return godog.ErrPending
}

func theConnectionToTheServerIsReestablished() error {
	return godog.ErrPending
}

func theServerReceivedTheMessageELeventPrefix() error {
	return godog.ErrPending
}

func theServerReceivedTheMessageEEVTtestSyetAnotherValue(arg1 int) error {
	return godog.ErrPending
}

func theClientUnlistensToEventsMatching(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsEULeventPrefix() error {
	return godog.ErrPending
}

func theServerSendsTheMessageESPeventPrefixeventPrefixfoundAMatch() error {
	return godog.ErrPending
}

func theClientWillBeNotifiedOfNewEventMatch(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageESReventPrefixeventPrefixfoundAMatch() error {
	return godog.ErrPending
}

func theClientWillBeNotifiedOfEventMatchRemoval(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageEAULeventPrefix() error {
	return godog.ErrPending
}

func theClientUnsubscribesFromAnEventNamed(arg1 string) error {
	return godog.ErrPending
}

func theServerReceivedTheMessageEUStest(arg1 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageEEVTtestSsomeValue(arg1 int) error {
	return godog.ErrPending
}

func theClientReceivedTheEventWithData(arg1, arg2 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageEEVTtestSanotherValue(arg1 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageEAUStest(arg1 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageIOnlyHaveOnePart() error {
	return godog.ErrPending
}

func theServerSendsTheMessageBR() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRXXX() error {
	return godog.ErrPending
}

func theClientSubscribesToPresenceEvents() error {
	return godog.ErrPending
}

func theServerReceivedTheMessageUSS() error {
	return godog.ErrPending
}

func theServerSendsTheMessageUASU() error {
	return godog.ErrPending
}

func theClientQueriesForConnectedClients() error {
	return godog.ErrPending
}

func theServerReceivedTheMessageUQQ() error {
	return godog.ErrPending
}

func theServerSendsTheMessageUQ() error {
	return godog.ErrPending
}

func theClientIsNotifiedThatNoClientsAreConnected() error {
	return godog.ErrPending
}

func theServerSendsTheMessageUQHomerMargeBart() error {
	return godog.ErrPending
}

func theClientIsNotifiedThatClientsAreConnected(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageUPNJHomer() error {
	return godog.ErrPending
}

func theClientIsNotifiedThatClientLoggedIn(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageUPNLBart() error {
	return godog.ErrPending
}

func theClientIsNotifiedThatClientLoggedOut(arg1 string) error {
	return godog.ErrPending
}

func theClientUnsubscribesToPresenceEvents() error {
	return godog.ErrPending
}

func theServerReceivedTheMessageUUSUS() error {
	return godog.ErrPending
}

func theServerSendsTheMessageUAUSU() error {
	return godog.ErrPending
}

func theClientIsNotNotifiedThatClientLoggedIn(arg1 string) error {
	return godog.ErrPending
}

func theClientCreatesARecordNamed(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRASmergeRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRRmergeRecordValue(arg1 int, arg2 string, arg3 int) error {
	return godog.ErrPending
}

func theClientSetsTheRecordKeyTo(arg1, arg2 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageREVERSIONEXISTSmergeRecordValue(arg1 int, arg2 string, arg3 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRUmergeRecordValue(arg1 int, arg2 string, arg3 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRUmergeRecordValue(arg1 int, arg2 string, arg3 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRCRconnectionRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRASconnectionRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRRconnectionRecordJohnRufflesDog(arg1 int, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theClientListensToARecordMatching(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRLrecordPrefix() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRALrecordPrefix() error {
	return godog.ErrPending
}

func theClientSetsTheRecordPetsNameTo(arg1 string, arg2 int, arg3 string) error {
	return godog.ErrPending
}

func theServerReceivedTheMessageRCRconnectionRecord() error {
	return godog.ErrPending
}

func theServerReceivedTheMessageRLrecordPrefix() error {
	return godog.ErrPending
}

func theServerReceivedTheMessageRPconnectionRecordPetsNameSMax(arg1, arg2 int) error {
	return godog.ErrPending
}

func theClientChecksIfTheServerHasTheRecord(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRHexistingRecordT() error {
	return godog.ErrPending
}

func theClientIsToldTheRecordExists(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRHnonExistentRecordF() error {
	return godog.ErrPending
}

func theClientIsToldTheRecordDoesntExist(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRAShasRecord() error {
	return godog.ErrPending
}

func theServerDidntReceiveTheMessageRHhasRecord() error {
	return godog.ErrPending
}

func theClientUnlistensToARecordMatching(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRULrecordPrefix() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRSPrecordPrefixrecordPrefixfoundAMatch() error {
	return godog.ErrPending
}

func theClientWillBeNotifiedOfNewRecordMatch(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRSRrecordPrefixrecordPrefixfoundAMatch() error {
	return godog.ErrPending
}

func theClientWillBeNotifiedOfRecordMatchRemoval(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRAULrecordPrefix() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRASdoubleRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRRdoubleRecordJohnRufflesDog(arg1 int, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRCRdoubleRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRASsnapshotRecord() error {
	return godog.ErrPending
}

func theClientRequestsASnapshotForTheRecord(arg1 string) error {
	return godog.ErrPending
}

func theClientHasNoResponseForTheSnapshotOfRecord(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRESNsnapshotRecordRECORDNOTFOUND() error {
	return godog.ErrPending
}

func theClientIsToldTheRecordEncounteredAnErrorRetrievingSnapshot(arg1 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRRsnapshotRecordJohn(arg1 int, arg2 string) error {
	return godog.ErrPending
}

func theClientIsProvidedTheSnapshotForRecordWithDataJohn(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRCRsubscribeRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRASsubscribeRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRRsubscribeRecordSmithRuffusDog(arg1 int, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theClientSubscribesToTheEntireRecordChanges(arg1 string) error {
	return godog.ErrPending
}

func theClientWillNotBeNotifiedOfTheRecordChange() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRUsubscribeRecordSmithRuffusDog(arg1 int, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theClientWillBeNotifiedOfTheRecordChange() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRPsubscribeRecordPetsNameSRuffusTheSecond(arg1, arg2 int) error {
	return godog.ErrPending
}

func theClientWillBeNotifiedOfThePartialRecordChange() error {
	return godog.ErrPending
}

func theClientUnsubscribesToTheEntireRecordChanges(arg1 string) error {
	return godog.ErrPending
}

func theClientSubscribesToForTheRecord(arg1, arg2 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRPsubscribeRecordNameSJohnSmith(arg1 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRUsubscribeRecordJohnSmithRuffusDog(arg1 int, arg2, arg3 string, arg4 int, arg5, arg6, arg7, arg8 string, arg9 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRPsubscribeRecordPetsAgeN(arg1, arg2, arg3 int) error {
	return godog.ErrPending
}

func theClientWillBeNotifiedOfTheSecondRecordChange() error {
	return godog.ErrPending
}

func theClientUnsubscribesToForTheRecord(arg1, arg2 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRASunhappyRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRRunhappyRecord(arg1 int, arg2, arg3 string) error {
	return godog.ErrPending
}

func theClientSetsTheRecordTo(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRUunhappyRecord(arg1 int, arg2, arg3 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRECACHERETRIEVALTIMEOUTunhappyRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRESTORAGERETRIEVALTIMEOUTunhappyRecord() error {
	return godog.ErrPending
}

func theClientDiscardsTheRecordNamed(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRUSunhappyRecord() error {
	return godog.ErrPending
}

func theClientDeletesTheRecordNamed(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRDunhappyRecord() error {
	return godog.ErrPending
}

func theClientRequiresWriteAcknowledgementOnRecord(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRCRhappyRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRAShappyRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRRhappyRecordJohnRufflesDog(arg1 int, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theClientRecordDataIsJohnRufflesDog(arg1, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRPhappyRecordPetsNameSMaxTrue(arg1, arg2 int, arg3 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRWAhappyRecordL(arg1 int) error {
	return godog.ErrPending
}

func theClientIsNotifiedThatTheRecordWasWrittenWithoutError(arg1 string) error {
	return godog.ErrPending
}

func theClientSetsTheRecordToSomeValue(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRUhappyRecordSomeValueTrue(arg1 int, arg2, arg3 string) error {
	return godog.ErrPending
}

func theClientRecordDataIsSomeValue(arg1, arg2 string) error {
	return godog.ErrPending
}

func theClientSetsTheRecordToNewErrorData(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRUhappyRecordNewErrorDataTrue(arg1 int, arg2, arg3 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRWAhappyRecordSErrorWritingRecordToStorage(arg1 int) error {
	return godog.ErrPending
}

func theClientIsNotifiedThatTheRecordWasWrittenWithError(arg1, arg2 string) error {
	return godog.ErrPending
}

func theClientSetsTheRecordValidDataTo(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRPhappyRecordValidDataSdifferentDataTrue(arg1 int, arg2 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRWAhappyRecordSErrorWritingRecordToCache(arg1 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRPhappyRecordValidDataSsomeDataTrue(arg1 int, arg2 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageREVERSIONEXISTShappyRecordNewErrorDataTrue(arg1 int, arg2, arg3 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRPhappyRecordPetsAgeN(arg1, arg2, arg3 int) error {
	return godog.ErrPending
}

func theServerSendsTheMessageRUhappyRecordSmithRuffusDog(arg1 int, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theClientRecordDataIsSmithRuffusDog(arg1, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRPhappyRecordPetsNameSMax(arg1, arg2 int) error {
	return godog.ErrPending
}

func theClientSetsTheRecordToSmithRuffusDog(arg1, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRUhappyRecordSmithRuffusDog(arg1 int, arg2, arg3, arg4, arg5, arg6 string, arg7 int) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRUShappyRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRAUShappyRecord() error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsRDhappyRecord() error {
	return godog.ErrPending
}

func theServerSendsTheMessageRADhappyRecord() error {
	return godog.ErrPending
}

func theClientProvidesARPCCalled(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPStoUppercase() error {
	return godog.ErrPending
}

func theServerSendsTheMessagePAStoUppercase() error {
	return godog.ErrPending
}

func theServerReceivedTheMessagePStoUppercase() error {
	return godog.ErrPending
}

func theServerSendsTheMessagePREQtoUppercaseUIDSabc() error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPAREQtoUppercaseUID() error {
	return godog.ErrPending
}

func theClientRecievesARequestForARPCCalledWithData(arg1, arg2 string) error {
	return godog.ErrPending
}

func theClientRespondsToTheRPCWithData(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPREStoUppercaseUIDSABC() error {
	return godog.ErrPending
}

func theClientRespondsToTheRPCWithTheError(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPEAnErrorOccuredtoUppercaseUID() error {
	return godog.ErrPending
}

func theClientRejectsTheRPC(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPREJtoUppercaseUID() error {
	return godog.ErrPending
}

func theServerSendsTheMessagePREQunSupportedUIDSabc() error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPREJunSupportedUID() error {
	return godog.ErrPending
}

func theClientStopsProvidingARPCCalled(arg1 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPUStoUppercase() error {
	return godog.ErrPending
}

func theServerSendsTheMessagePAUStoUppercase() error {
	return godog.ErrPending
}

func theClientRequestsRPCWithData(arg1, arg2 string) error {
	return godog.ErrPending
}

func theLastMessageTheServerRecievedIsPREQtoUppercaseUIDSabc() error {
	return godog.ErrPending
}

func theServerSendsTheMessagePAREQtoUppercaseUID() error {
	return godog.ErrPending
}

func theServerSendsTheMessagePREStoUppercaseUIDSABC() error {
	return godog.ErrPending
}

func theClientRecievesASuccessfulRPCCallbackForWithData(arg1, arg2 string) error {
	return godog.ErrPending
}

func theServerSendsTheMessagePAREQtoUpperCaseUID() error {
	return godog.ErrPending
}

func theServerSendsTheMessagePERPCErrorMessagetoUppercaseUID() error {
	return godog.ErrPending
}

func theClientRecievesAnErrorRPCCallbackForWithTheMessage(arg1, arg2 string) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.AfterScenario(afterScenario)

	s.Step(`^the test server is ready$`, theTestServerIsReady)
	s.Step(`^the server has (\d+) active connections$`, theServerHasActiveConnections)
	s.Step(`^the client is initialised$`, theClientIsInitialised)
	s.Step(`^the clients connection state is "([^"]*)"$`, theClientsConnectionStateIs)
	s.Step(`^the client is initialised with a small heartbeat interval$`, theClientIsInitialisedWithASmallHeartbeatInterval)
	s.Step(`^the server sends the message C\|A\+$`, theServerSendsTheMessage("C", "A"))
	s.Step(`^the client logs in with username "([^"]*)" and password "([^"]*)"$`, theClientLogsInWithUsernameAndPassword)
	s.Step(`^the server sends the message A\|A\+$`, theServerSendsTheMessage("A", "A"))
	s.Step(`^the server sends the message C\|PI\+$`, theServerSendsTheMessage("C", "PI"))
	s.Step(`^the server received the message C\|PO\+$`, theServerReceivedTheMessage("C", "PO"))
	s.Step(`^two seconds later$`, someMsLater(2000))
	s.Step(`^the client throws a "([^"]*)" error with message "([^"]*)"$`, theClientThrowsAnErrorWithMessage)
	s.Step(`^the second test server is ready$`, theSecondTestServerIsReady)
	s.Step(`^the second server has (\d+) active connections$`, theSecondServerHasActiveConnections)
	s.Step(`^the server sends the message C\|CH\+$`, theServerSendsTheMessageCCH)
	s.Step(`^the last message the server recieved is C\|CHR\|<FIRST_SERVER_URL>\+$`, theLastMessageTheServerRecievedIsCCHRFIRSTSERVERURL)
	s.Step(`^the server sends the message C\|REJ\+$`, theServerSendsTheMessageCREJ)
	s.Step(`^the server has received (\d+) messages$`, theServerHasReceivedMessages)
	s.Step(`^the server sends the message C\|RED\|<SECOND_SERVER_URL>\+$`, theServerSendsTheMessageCREDSECONDSERVERURL)
	s.Step(`^some time passes$`, someTimePasses)
	s.Step(`^the client is on the second server$`, theClientIsOnTheSecondServer)
	s.Step(`^the last message the server recieved is A\|REQ\|{"([^"]*)":"XXX","([^"]*)":"YYY"}\+$`, theLastMessageTheServerRecievedIsAREQXXXYYY)
	s.Step(`^the last login was successful$`, theLastLoginWasSuccessful)
	s.Step(`^the server sends the message A\|E\|INVALID_AUTH_DATA\|Sinvalid authentication data\+$`, theServerSendsTheMessageAEINVALIDAUTHDATASinvalidAuthenticationData)
	s.Step(`^the last login failed with error message "([^"]*)"$`, theLastLoginFailedWithErrorMessage)
	s.Step(`^the server sends the message A\|E\|TOO_MANY_AUTH_ATTEMPTS\|Stoo many authentication attempts\+$`, theServerSendsTheMessageAETOOMANYAUTHATTEMPTSStooManyAuthenticationAttempts)
	s.Step(`^the server resets its message count$`, theServerResetsItsMessageCount)
	s.Step(`^the client subscribes to an event named "([^"]*)"$`, theClientSubscribesToAnEventNamed)
	s.Step(`^the server sends the message E\|A\|S\|test(\d+)\+$`, theServerSendsTheMessageEAStest)
	s.Step(`^the server received the message E\|S\|test(\d+)\+$`, theServerReceivedTheMessageEStest)
	s.Step(`^the last message the server recieved is E\|S\|test(\d+)\+$`, theLastMessageTheServerRecievedIsEStest)
	s.Step(`^the client listens to events matching "([^"]*)"$`, theClientListensToEventsMatching)
	s.Step(`^the last message the server recieved is E\|L\|eventPrefix\/\.\*\+$`, theLastMessageTheServerRecievedIsELeventPrefix)
	s.Step(`^the server sends the message E\|A\|L\|eventPrefix\/\.\*\+$`, theServerSendsTheMessageEALeventPrefix)
	s.Step(`^the connection to the server is lost$`, theConnectionToTheServerIsLost)
	s.Step(`^the client publishes an event named "([^"]*)" with data "([^"]*)"$`, theClientPublishesAnEventNamedWithData)
	s.Step(`^the server did not recieve any messages$`, theServerDidNotRecieveAnyMessages)
	s.Step(`^the connection to the server is reestablished$`, theConnectionToTheServerIsReestablished)
	s.Step(`^the server received the message E\|L\|eventPrefix\/\.\*\+$`, theServerReceivedTheMessageELeventPrefix)
	s.Step(`^the server received the message E\|EVT\|test(\d+)\|SyetAnotherValue\+$`, theServerReceivedTheMessageEEVTtestSyetAnotherValue)
	s.Step(`^the client unlistens to events matching "([^"]*)"$`, theClientUnlistensToEventsMatching)
	s.Step(`^the last message the server recieved is E\|UL\|eventPrefix\/\.\*\+$`, theLastMessageTheServerRecievedIsEULeventPrefix)
	s.Step(`^the server sends the message E\|SP\|eventPrefix\/\.\*\|eventPrefix\/foundAMatch\+$`, theServerSendsTheMessageESPeventPrefixeventPrefixfoundAMatch)
	s.Step(`^the client will be notified of new event match "([^"]*)"$`, theClientWillBeNotifiedOfNewEventMatch)
	s.Step(`^the server sends the message E\|SR\|eventPrefix\/\.\*\|eventPrefix\/foundAMatch\+$`, theServerSendsTheMessageESReventPrefixeventPrefixfoundAMatch)
	s.Step(`^the client will be notified of event match removal "([^"]*)"$`, theClientWillBeNotifiedOfEventMatchRemoval)
	s.Step(`^the server sends the message E\|A\|UL\|eventPrefix\/\.\*\+$`, theServerSendsTheMessageEAULeventPrefix)
	s.Step(`^the client unsubscribes from an event named "([^"]*)"$`, theClientUnsubscribesFromAnEventNamed)
	s.Step(`^the server received the message E\|US\|test(\d+)\+$`, theServerReceivedTheMessageEUStest)
	s.Step(`^the server sends the message E\|EVT\|test(\d+)\|SsomeValue\+$`, theServerSendsTheMessageEEVTtestSsomeValue)
	s.Step(`^the client received the event "([^"]*)" with data "([^"]*)"$`, theClientReceivedTheEventWithData)
	s.Step(`^the server sends the message E\|EVT\|test(\d+)\|SanotherValue\+$`, theServerSendsTheMessageEEVTtestSanotherValue)
	s.Step(`^the server sends the message E\|A\|US\|test(\d+)\+$`, theServerSendsTheMessageEAUStest)
	s.Step(`^the server sends the message I only have one part\+$`, theServerSendsTheMessageIOnlyHaveOnePart)
	s.Step(`^the server sends the message B\|R\+$`, theServerSendsTheMessageBR)
	s.Step(`^the server sends the message R\|XXX\+$`, theServerSendsTheMessageRXXX)
	s.Step(`^the client subscribes to presence events$`, theClientSubscribesToPresenceEvents)
	s.Step(`^the server received the message U\|S\|S\+$`, theServerReceivedTheMessageUSS)
	s.Step(`^the server sends the message U\|A\|S\|U\+$`, theServerSendsTheMessageUASU)
	s.Step(`^the client queries for connected clients$`, theClientQueriesForConnectedClients)
	s.Step(`^the server received the message U\|Q\|Q\+$`, theServerReceivedTheMessageUQQ)
	s.Step(`^the server sends the message U\|Q\+$`, theServerSendsTheMessageUQ)
	s.Step(`^the client is notified that no clients are connected$`, theClientIsNotifiedThatNoClientsAreConnected)
	s.Step(`^the server sends the message U\|Q\|Homer\|Marge\|Bart\+$`, theServerSendsTheMessageUQHomerMargeBart)
	s.Step(`^the client is notified that clients "([^"]*)" are connected$`, theClientIsNotifiedThatClientsAreConnected)
	s.Step(`^the server sends the message U\|PNJ\|Homer\+$`, theServerSendsTheMessageUPNJHomer)
	s.Step(`^the client is notified that client "([^"]*)" logged in$`, theClientIsNotifiedThatClientLoggedIn)
	s.Step(`^the server sends the message U\|PNL\|Bart\+$`, theServerSendsTheMessageUPNLBart)
	s.Step(`^the client is notified that client "([^"]*)" logged out$`, theClientIsNotifiedThatClientLoggedOut)
	s.Step(`^the client unsubscribes to presence events$`, theClientUnsubscribesToPresenceEvents)
	s.Step(`^the server received the message U\|US\|US\+$`, theServerReceivedTheMessageUUSUS)
	s.Step(`^the server sends the message U\|A\|US\|U\+$`, theServerSendsTheMessageUAUSU)
	s.Step(`^the client is not notified that client "([^"]*)" logged in$`, theClientIsNotNotifiedThatClientLoggedIn)
	s.Step(`^the client creates a record named "([^"]*)"$`, theClientCreatesARecordNamed)
	s.Step(`^the server sends the message R\|A\|S\|mergeRecord\+$`, theServerSendsTheMessageRASmergeRecord)
	s.Step(`^the server sends the message R\|R\|mergeRecord\|(\d+)\|{"([^"]*)":"value(\d+)"}\+$`, theServerSendsTheMessageRRmergeRecordValue)
	s.Step(`^the client sets the record "([^"]*)" "key" to "([^"]*)"$`, theClientSetsTheRecordKeyTo)
	s.Step(`^the server sends the message R\|E\|VERSION_EXISTS\|mergeRecord\|(\d+)\|{"([^"]*)":"value(\d+)"}\+$`, theServerSendsTheMessageREVERSIONEXISTSmergeRecordValue)
	s.Step(`^the last message the server recieved is R\|U\|mergeRecord\|(\d+)\|{"([^"]*)":"value(\d+)"}\+$`, theLastMessageTheServerRecievedIsRUmergeRecordValue)
	s.Step(`^the server sends the message R\|U\|mergeRecord\|(\d+)\|{"([^"]*)":"value(\d+)"}\+$`, theServerSendsTheMessageRUmergeRecordValue)
	s.Step(`^the last message the server recieved is R\|CR\|connectionRecord\+$`, theLastMessageTheServerRecievedIsRCRconnectionRecord)
	s.Step(`^the server sends the message R\|A\|S\|connectionRecord\+$`, theServerSendsTheMessageRASconnectionRecord)
	s.Step(`^the server sends the message R\|R\|connectionRecord\|(\d+)\|{"([^"]*)":"John", "([^"]*)": \[{"([^"]*)":"Ruffles", "([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theServerSendsTheMessageRRconnectionRecordJohnRufflesDog)
	s.Step(`^the client listens to a record matching "([^"]*)"$`, theClientListensToARecordMatching)
	s.Step(`^the last message the server recieved is R\|L\|recordPrefix\/\.\*\+$`, theLastMessageTheServerRecievedIsRLrecordPrefix)
	s.Step(`^the server sends the message R\|A\|L\|recordPrefix\/\.\*\+$`, theServerSendsTheMessageRALrecordPrefix)
	s.Step(`^the client sets the record "([^"]*)" "pets\.(\d+)\.name" to "([^"]*)"$`, theClientSetsTheRecordPetsNameTo)
	s.Step(`^the server received the message R\|CR\|connectionRecord\+$`, theServerReceivedTheMessageRCRconnectionRecord)
	s.Step(`^the server received the message R\|L\|recordPrefix\/\.\*\+$`, theServerReceivedTheMessageRLrecordPrefix)
	s.Step(`^the server received the message R\|P\|connectionRecord\|(\d+)\|pets\.(\d+)\.name\|SMax\+$`, theServerReceivedTheMessageRPconnectionRecordPetsNameSMax)
	s.Step(`^the client checks if the server has the record "([^"]*)"$`, theClientChecksIfTheServerHasTheRecord)
	s.Step(`^the server sends the message R\|H\|existingRecord\|T\|\+$`, theServerSendsTheMessageRHexistingRecordT)
	s.Step(`^the client is told the record "([^"]*)" exists$`, theClientIsToldTheRecordExists)
	s.Step(`^the server sends the message R\|H\|nonExistentRecord\|F\|\+$`, theServerSendsTheMessageRHnonExistentRecordF)
	s.Step(`^the client is told the record "([^"]*)" doesn\'t exist$`, theClientIsToldTheRecordDoesntExist)
	s.Step(`^the server sends the message R\|A\|S\|hasRecord\+$`, theServerSendsTheMessageRAShasRecord)
	s.Step(`^the server didn\'t receive the message R\|H\|hasRecord\+$`, theServerDidntReceiveTheMessageRHhasRecord)
	s.Step(`^the client unlistens to a record matching "([^"]*)"$`, theClientUnlistensToARecordMatching)
	s.Step(`^the last message the server recieved is R\|UL\|recordPrefix\/\.\*\+$`, theLastMessageTheServerRecievedIsRULrecordPrefix)
	s.Step(`^the server sends the message R\|SP\|recordPrefix\/\.\*\|recordPrefix\/foundAMatch\+$`, theServerSendsTheMessageRSPrecordPrefixrecordPrefixfoundAMatch)
	s.Step(`^the client will be notified of new record match "([^"]*)"$`, theClientWillBeNotifiedOfNewRecordMatch)
	s.Step(`^the server sends the message R\|SR\|recordPrefix\/\.\*\|recordPrefix\/foundAMatch\+$`, theServerSendsTheMessageRSRrecordPrefixrecordPrefixfoundAMatch)
	s.Step(`^the client will be notified of record match removal "([^"]*)"$`, theClientWillBeNotifiedOfRecordMatchRemoval)
	s.Step(`^the server sends the message R\|A\|UL\|recordPrefix\/\.\*\+$`, theServerSendsTheMessageRAULrecordPrefix)
	s.Step(`^the server sends the message R\|A\|S\|doubleRecord\+$`, theServerSendsTheMessageRASdoubleRecord)
	s.Step(`^the server sends the message R\|R\|doubleRecord\|(\d+)\|{"([^"]*)":"John", "([^"]*)": \[{"([^"]*)":"Ruffles", "([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theServerSendsTheMessageRRdoubleRecordJohnRufflesDog)
	s.Step(`^the last message the server recieved is R\|CR\|doubleRecord\+$`, theLastMessageTheServerRecievedIsRCRdoubleRecord)
	s.Step(`^the server sends the message R\|A\|S\|snapshotRecord\+$`, theServerSendsTheMessageRASsnapshotRecord)
	s.Step(`^the client requests a snapshot for the record "([^"]*)"$`, theClientRequestsASnapshotForTheRecord)
	s.Step(`^the client has no response for the snapshot of record "([^"]*)"$`, theClientHasNoResponseForTheSnapshotOfRecord)
	s.Step(`^the server sends the message R\|E\|SN\|snapshotRecord\|RECORD_NOT_FOUND\+$`, theServerSendsTheMessageRESNsnapshotRecordRECORDNOTFOUND)
	s.Step(`^the client is told the record "([^"]*)" encountered an error retrieving snapshot$`, theClientIsToldTheRecordEncounteredAnErrorRetrievingSnapshot)
	s.Step(`^the server sends the message R\|R\|snapshotRecord\|(\d+)\|{"([^"]*)":"John"}\+$`, theServerSendsTheMessageRRsnapshotRecordJohn)
	s.Step(`^the client is provided the snapshot for record "([^"]*)" with data "{"([^"]*)":"John"}"$`, theClientIsProvidedTheSnapshotForRecordWithDataJohn)
	s.Step(`^the last message the server recieved is R\|CR\|subscribeRecord\+$`, theLastMessageTheServerRecievedIsRCRsubscribeRecord)
	s.Step(`^the server sends the message R\|A\|S\|subscribeRecord\+$`, theServerSendsTheMessageRASsubscribeRecord)
	s.Step(`^the server sends the message R\|R\|subscribeRecord\|(\d+)\|{"([^"]*)":"Smith","([^"]*)":\[{"([^"]*)":"Ruffus","([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theServerSendsTheMessageRRsubscribeRecordSmithRuffusDog)
	s.Step(`^the client subscribes to the entire record "([^"]*)" changes$`, theClientSubscribesToTheEntireRecordChanges)
	s.Step(`^the client will not be notified of the record change$`, theClientWillNotBeNotifiedOfTheRecordChange)
	s.Step(`^the server sends the message R\|U\|subscribeRecord\|(\d+)\|{"([^"]*)":"Smith","([^"]*)":\[{"([^"]*)":"Ruffus","([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theServerSendsTheMessageRUsubscribeRecordSmithRuffusDog)
	s.Step(`^the client will be notified of the record change$`, theClientWillBeNotifiedOfTheRecordChange)
	s.Step(`^the server sends the message R\|P\|subscribeRecord\|(\d+)\|pets\.(\d+)\.name\|SRuffusTheSecond\+$`, theServerSendsTheMessageRPsubscribeRecordPetsNameSRuffusTheSecond)
	s.Step(`^the client will be notified of the partial record change$`, theClientWillBeNotifiedOfThePartialRecordChange)
	s.Step(`^the client unsubscribes to the entire record "([^"]*)" changes$`, theClientUnsubscribesToTheEntireRecordChanges)
	s.Step(`^the client subscribes to "([^"]*)" for the record "([^"]*)"$`, theClientSubscribesToForTheRecord)
	s.Step(`^the server sends the message R\|P\|subscribeRecord\|(\d+)\|name\|SJohn Smith\+$`, theServerSendsTheMessageRPsubscribeRecordNameSJohnSmith)
	s.Step(`^the server sends the message R\|U\|subscribeRecord\|(\d+)\|{"([^"]*)":"John Smith", "([^"]*)": (\d+), "([^"]*)": \[{"([^"]*)":"Ruffus", "([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theServerSendsTheMessageRUsubscribeRecordJohnSmithRuffusDog)
	s.Step(`^the server sends the message R\|P\|subscribeRecord\|(\d+)\|pets\.(\d+)\.age\|N(\d+)\+$`, theServerSendsTheMessageRPsubscribeRecordPetsAgeN)
	s.Step(`^the client will be notified of the second record change$`, theClientWillBeNotifiedOfTheSecondRecordChange)
	s.Step(`^the client unsubscribes to "([^"]*)" for the record "([^"]*)"$`, theClientUnsubscribesToForTheRecord)
	s.Step(`^the server sends the message R\|A\|S\|unhappyRecord\+$`, theServerSendsTheMessageRASunhappyRecord)
	s.Step(`^the server sends the message R\|R\|unhappyRecord\|(\d+)\|{"([^"]*)":\["([^"]*)"\]}\+$`, theServerSendsTheMessageRRunhappyRecord)
	s.Step(`^the client sets the record "([^"]*)" to {"([^"]*)":\["([^"]*)"\]}$`, theClientSetsTheRecordTo)
	s.Step(`^the last message the server recieved is R\|U\|unhappyRecord\|(\d+)\|{"([^"]*)":\["([^"]*)"\]}\+$`, theLastMessageTheServerRecievedIsRUunhappyRecord)
	s.Step(`^the server sends the message R\|E\|CACHE_RETRIEVAL_TIMEOUT\|unhappyRecord\+$`, theServerSendsTheMessageRECACHERETRIEVALTIMEOUTunhappyRecord)
	s.Step(`^the server sends the message R\|E\|STORAGE_RETRIEVAL_TIMEOUT\|unhappyRecord\+$`, theServerSendsTheMessageRESTORAGERETRIEVALTIMEOUTunhappyRecord)
	s.Step(`^the client discards the record named "([^"]*)"$`, theClientDiscardsTheRecordNamed)
	s.Step(`^the last message the server recieved is R\|US\|unhappyRecord\+$`, theLastMessageTheServerRecievedIsRUSunhappyRecord)
	s.Step(`^the client deletes the record named "([^"]*)"$`, theClientDeletesTheRecordNamed)
	s.Step(`^the last message the server recieved is R\|D\|unhappyRecord\+$`, theLastMessageTheServerRecievedIsRDunhappyRecord)
	s.Step(`^the client requires write acknowledgement on record "([^"]*)"$`, theClientRequiresWriteAcknowledgementOnRecord)
	s.Step(`^the last message the server recieved is R\|CR\|happyRecord\+$`, theLastMessageTheServerRecievedIsRCRhappyRecord)
	s.Step(`^the server sends the message R\|A\|S\|happyRecord\+$`, theServerSendsTheMessageRAShappyRecord)
	s.Step(`^the server sends the message R\|R\|happyRecord\|(\d+)\|{"([^"]*)":"John", "([^"]*)": \[{"([^"]*)":"Ruffles", "([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theServerSendsTheMessageRRhappyRecordJohnRufflesDog)
	s.Step(`^the client record "([^"]*)" data is {"([^"]*)":"John", "([^"]*)": \[{"([^"]*)":"Ruffles", "([^"]*)":"dog","([^"]*)":(\d+)}\]}$`, theClientRecordDataIsJohnRufflesDog)
	s.Step(`^the last message the server recieved is R\|P\|happyRecord\|(\d+)\|pets\.(\d+)\.name\|SMax\|{"([^"]*)":true}\+$`, theLastMessageTheServerRecievedIsRPhappyRecordPetsNameSMaxTrue)
	s.Step(`^the server sends the message R\|WA\|happyRecord\|\[(\d+)\]\|L\+$`, theServerSendsTheMessageRWAhappyRecordL)
	s.Step(`^the client is notified that the record "([^"]*)" was written without error$`, theClientIsNotifiedThatTheRecordWasWrittenWithoutError)
	s.Step(`^the client sets the record "([^"]*)" to {"([^"]*)":"someValue"}$`, theClientSetsTheRecordToSomeValue)
	s.Step(`^the last message the server recieved is R\|U\|happyRecord\|(\d+)\|{"([^"]*)":"someValue"}\|{"([^"]*)":true}\+$`, theLastMessageTheServerRecievedIsRUhappyRecordSomeValueTrue)
	s.Step(`^the client record "([^"]*)" data is {"([^"]*)":"someValue"}$`, theClientRecordDataIsSomeValue)
	s.Step(`^the client sets the record "([^"]*)" to {"([^"]*)":"newErrorData"}$`, theClientSetsTheRecordToNewErrorData)
	s.Step(`^the last message the server recieved is R\|U\|happyRecord\|(\d+)\|{"([^"]*)":"newErrorData"}\|{"([^"]*)":true}\+$`, theLastMessageTheServerRecievedIsRUhappyRecordNewErrorDataTrue)
	s.Step(`^the server sends the message R\|WA\|happyRecord\|\[(\d+)\]\|SError writing record to storage\+$`, theServerSendsTheMessageRWAhappyRecordSErrorWritingRecordToStorage)
	s.Step(`^the client is notified that the record "([^"]*)" was written with error "([^"]*)"$`, theClientIsNotifiedThatTheRecordWasWrittenWithError)
	s.Step(`^the client sets the record "([^"]*)" "validData" to "([^"]*)"$`, theClientSetsTheRecordValidDataTo)
	s.Step(`^the last message the server recieved is R\|P\|happyRecord\|(\d+)\|validData\|SdifferentData\|{"([^"]*)":true}\+$`, theLastMessageTheServerRecievedIsRPhappyRecordValidDataSdifferentDataTrue)
	s.Step(`^the server sends the message R\|WA\|happyRecord\|\[(\d+)\]\|SError writing record to cache\+$`, theServerSendsTheMessageRWAhappyRecordSErrorWritingRecordToCache)
	s.Step(`^the last message the server recieved is R\|P\|happyRecord\|(\d+)\|validData\|SsomeData\|{"([^"]*)":true}\+$`, theLastMessageTheServerRecievedIsRPhappyRecordValidDataSsomeDataTrue)
	s.Step(`^the server sends the message R\|E\|VERSION_EXISTS\|happyRecord\|(\d+)\|{"([^"]*)":"newErrorData"}\|{"([^"]*)":true}\+$`, theServerSendsTheMessageREVERSIONEXISTShappyRecordNewErrorDataTrue)
	s.Step(`^the server sends the message R\|P\|happyRecord\|(\d+)\|pets\.(\d+)\.age\|N(\d+)\+$`, theServerSendsTheMessageRPhappyRecordPetsAgeN)
	s.Step(`^the server sends the message R\|U\|happyRecord\|(\d+)\|{"([^"]*)":"Smith", "([^"]*)": \[{"([^"]*)":"Ruffus", "([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theServerSendsTheMessageRUhappyRecordSmithRuffusDog)
	s.Step(`^the client record "([^"]*)" data is {"([^"]*)":"Smith", "([^"]*)": \[{"([^"]*)":"Ruffus", "([^"]*)":"dog","([^"]*)":(\d+)}\]}$`, theClientRecordDataIsSmithRuffusDog)
	s.Step(`^the last message the server recieved is R\|P\|happyRecord\|(\d+)\|pets\.(\d+)\.name\|SMax\+$`, theLastMessageTheServerRecievedIsRPhappyRecordPetsNameSMax)
	s.Step(`^the client sets the record "([^"]*)" to {"([^"]*)":"Smith","([^"]*)":\[{"([^"]*)":"Ruffus","([^"]*)":"dog","([^"]*)":(\d+)}\]}$`, theClientSetsTheRecordToSmithRuffusDog)
	s.Step(`^the last message the server recieved is R\|U\|happyRecord\|(\d+)\|{"([^"]*)":"Smith","([^"]*)":\[{"([^"]*)":"Ruffus","([^"]*)":"dog","([^"]*)":(\d+)}\]}\+$`, theLastMessageTheServerRecievedIsRUhappyRecordSmithRuffusDog)
	s.Step(`^the last message the server recieved is R\|US\|happyRecord\+$`, theLastMessageTheServerRecievedIsRUShappyRecord)
	s.Step(`^the server sends the message R\|A\|US\|happyRecord\+$`, theServerSendsTheMessageRAUShappyRecord)
	s.Step(`^the last message the server recieved is R\|D\|happyRecord\+$`, theLastMessageTheServerRecievedIsRDhappyRecord)
	s.Step(`^the server sends the message R\|A\|D\|happyRecord\+$`, theServerSendsTheMessageRADhappyRecord)
	s.Step(`^the client provides a RPC called "([^"]*)"$`, theClientProvidesARPCCalled)
	s.Step(`^the last message the server recieved is P\|S\|toUppercase\+$`, theLastMessageTheServerRecievedIsPStoUppercase)
	s.Step(`^the server sends the message P\|A\|S\|toUppercase\+$`, theServerSendsTheMessagePAStoUppercase)
	s.Step(`^the server received the message P\|S\|toUppercase\+$`, theServerReceivedTheMessagePStoUppercase)
	s.Step(`^the server sends the message P\|REQ\|toUppercase\|<UID>\|Sabc\+$`, theServerSendsTheMessagePREQtoUppercaseUIDSabc)
	s.Step(`^the last message the server recieved is P\|A\|REQ\|toUppercase\|<UID>\+$`, theLastMessageTheServerRecievedIsPAREQtoUppercaseUID)
	s.Step(`^the client recieves a request for a RPC called "([^"]*)" with data "([^"]*)"$`, theClientRecievesARequestForARPCCalledWithData)
	s.Step(`^the client responds to the RPC "([^"]*)" with data "([^"]*)"$`, theClientRespondsToTheRPCWithData)
	s.Step(`^the last message the server recieved is P\|RES\|toUppercase\|<UID>\|SABC\+$`, theLastMessageTheServerRecievedIsPREStoUppercaseUIDSABC)
	s.Step(`^the client responds to the RPC "([^"]*)" with the error "([^"]*)"$`, theClientRespondsToTheRPCWithTheError)
	s.Step(`^the last message the server recieved is P\|E\|An Error Occured\|toUppercase\|<UID>\+$`, theLastMessageTheServerRecievedIsPEAnErrorOccuredtoUppercaseUID)
	s.Step(`^the client rejects the RPC "([^"]*)"$`, theClientRejectsTheRPC)
	s.Step(`^the last message the server recieved is P\|REJ\|toUppercase\|<UID>\+$`, theLastMessageTheServerRecievedIsPREJtoUppercaseUID)
	s.Step(`^the server sends the message P\|REQ\|unSupported\|<UID>\|Sabc\+$`, theServerSendsTheMessagePREQunSupportedUIDSabc)
	s.Step(`^the last message the server recieved is P\|REJ\|unSupported\|<UID>\+$`, theLastMessageTheServerRecievedIsPREJunSupportedUID)
	s.Step(`^the client stops providing a RPC called "([^"]*)"$`, theClientStopsProvidingARPCCalled)
	s.Step(`^the last message the server recieved is P\|US\|toUppercase\+$`, theLastMessageTheServerRecievedIsPUStoUppercase)
	s.Step(`^the server sends the message P\|A\|US\|toUppercase\+$`, theServerSendsTheMessagePAUStoUppercase)
	s.Step(`^the client requests RPC "([^"]*)" with data "([^"]*)"$`, theClientRequestsRPCWithData)
	s.Step(`^the last message the server recieved is P\|REQ\|toUppercase\|<UID>\|Sabc\+$`, theLastMessageTheServerRecievedIsPREQtoUppercaseUIDSabc)
	s.Step(`^the server sends the message P\|A\|REQ\|toUppercase\|<UID>\+$`, theServerSendsTheMessagePAREQtoUppercaseUID)
	s.Step(`^the server sends the message P\|RES\|toUppercase\|<UID>\|SABC\+$`, theServerSendsTheMessagePREStoUppercaseUIDSABC)
	s.Step(`^the client recieves a successful RPC callback for "([^"]*)" with data "([^"]*)"$`, theClientRecievesASuccessfulRPCCallbackForWithData)
	s.Step(`^the server sends the message P\|A\|REQ\|toUpperCase\|<UID>\+$`, theServerSendsTheMessagePAREQtoUpperCaseUID)
	s.Step(`^the server sends the message P\|E\|RPC Error Message\|toUppercase\|<UID>\+$`, theServerSendsTheMessagePERPCErrorMessagetoUppercaseUID)
	s.Step(`^the client recieves an error RPC callback for "([^"]*)" with the message "([^"]*)"$`, theClientRecievesAnErrorRPCCallbackForWithTheMessage)
}
