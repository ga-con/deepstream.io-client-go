// deepstream.io-client-go
// https://github.com/ga-con/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package interfaces

//MessageSeparator invisible message separator
const MessageSeparator = string(byte(30))

//MessagePartSeparator invisible message separator
const MessagePartSeparator = string(byte(31))

//SourceMessageConnector identifies source message connector key
const SourceMessageConnector = "SOURCE_MESSAGE_CONNECTOR"

//Log Level

//LogLevelDebug identifies log level
const LogLevelDebug = 0

//LogLevelInfo identifies log level
const LogLevelInfo = 1

//LogLevelWarn identifies log level
const LogLevelWarn = 2

//LogLevelError identifies log level
const LogLevelError = 2

//LogLevelOff identifies log level
const LogLevelOff = 100

//Server State

//ServerStateStarting indicates the server is starting
const ServerStateStarting = "starting"

//ServerStateInitialized indicates the server is initialized
const ServerStateInitialized = "initialized"

//ServerStateRunning indicates the server is initialized
const ServerStateRunning = "is-running"

//ServerStateClosing indicates the server is closing
const ServerStateClosing = "closing"

//ServerStateClosed indicates the server is closed
const ServerStateClosed = "closed"

//Connection State

//ConnectionState indicates the possible connection states
type ConnectionState string

//ConnectionStateClosed indicates the connection has been closed
const ConnectionStateClosed ConnectionState = "CLOSED"

//ConnectionStateAwaitingConnection indicates the connection is waiting for a connection
const ConnectionStateAwaitingConnection ConnectionState = "AWAITING_CONNECTION"

//ConnectionStateChallenging indicates the connection is being challenged for credentials
const ConnectionStateChallenging ConnectionState = "CHALLENGING"

//ConnectionStateAwaitingAuthentication indicates the connection is waiting for authentication
const ConnectionStateAwaitingAuthentication ConnectionState = "AWAITING_AUTHENTICATION"

//ConnectionStateAuthenticating indicates the connection is authenticating with the server
const ConnectionStateAuthenticating ConnectionState = "AUTHENTICATING"

//ConnectionStateOpen indicates the connection is open
const ConnectionStateOpen ConnectionState = "OPEN"

//ConnectionStateError indicates the connection has errored
const ConnectionStateError ConnectionState = "ERROR"

//ConnectionStateReconnecting indicates the connection is reconnecting
const ConnectionStateReconnecting ConnectionState = "RECONNECTING"

//Event

//name	value	server	client
//EVENT.TRIGGER_EVENT	TRIGGER_EVENT	✔
//EVENT.INCOMING_CONNECTION	INCOMING_CONNECTION	✔
//EVENT.INFO	INFO	✔
//EVENT.SUBSCRIBE	SUBSCRIBE	✔
//EVENT.UNSUBSCRIBE	UNSUBSCRIBE	✔
//EVENT.RECORD_DELETION	RECORD_DELETION	✔
//EVENT.INVALID_AUTH_MSG	INVALID_AUTH_MSG	✔
//EVENT.INVALID_AUTH_DATA	INVALID_AUTH_DATA	✔
//EVENT.AUTH_ATTEMPT	AUTH_ATTEMPT	✔
//EVENT.AUTH_ERROR	AUTH_ERROR	✔
//EVENT.TOO_MANY_AUTH_ATTEMPTS	TOO_MANY_AUTH_ATTEMPTS	✔	✔
//EVENT.AUTH_SUCCESSFUL	AUTH_SUCCESSFUL	✔
//EVENT.NOT_AUTHENTICATED	NOT_AUTHENTICATED		✔
//EVENT.CONNECTION_ERROR	CONNECTION_ERROR	✔	✔
//EVENT.MESSAGE_PERMISSION_ERROR	MESSAGE_PERMISSION_ERROR	✔	✔
//EVENT.MESSAGE_PARSE_ERROR	MESSAGE_PARSE_ERROR	✔	✔
//EVENT.MAXIMUM_MESSAGE_SIZE_EXCEEDED	MAXIMUM_MESSAGE_SIZE_EXCEEDED	✔
//EVENT.MESSAGE_DENIED	MESSAGE_DENIED	✔	✔
//EVENT.INVALID_MESSAGE_DATA	INVALID_MESSAGE_DATA	✔
//EVENT.UNKNOWN_TOPIC	UNKNOWN_TOPIC	✔
//EVENT.UNKNOWN_ACTION	UNKNOWN_ACTION	✔
//EVENT.MULTIPLE_SUBSCRIPTIONS	MULTIPLE_SUBSCRIPTIONS	✔
//EVENT.NOT_SUBSCRIBED	NOT_SUBSCRIBED	✔
//EVENT.LISTENER_EXISTS	LISTENER_EXISTS		✔
//EVENT.NOT_LISTENING	NOT_LISTENING		✔
//EVENT.IS_CLOSED	IS_CLOSED		✔
//EVENT.ACK_TIMEOUT	ACK_TIMEOUT	✔	✔
//EVENT.RESPONSE_TIMEOUT	RESPONSE_TIMEOUT	✔	✔
//EVENT.DELETE_TIMEOUT	DELETE_TIMEOUT		✔
//EVENT.UNSOLICITED_MESSAGE	UNSOLICITED_MESSAGE		✔
//EVENT.MULTIPLE_ACK	MULTIPLE_ACK	✔
//EVENT.MULTIPLE_RESPONSE	MULTIPLE_RESPONSE	✔
//EVENT.NO_RPC_PROVIDER	NO_RPC_PROVIDER	✔
//EVENT.RECORD_LOAD_ERROR	RECORD_LOAD_ERROR	✔
//EVENT.RECORD_CREATE_ERROR	RECORD_CREATE_ERROR	✔
//EVENT.RECORD_UPDATE_ERROR	RECORD_UPDATE_ERROR	✔
//EVENT.RECORD_DELETE_ERROR	RECORD_DELETE_ERROR	✔
//EVENT.RECORD_SNAPSHOT_ERROR	RECORD_SNAPSHOT_ERROR	✔
//EVENT.RECORD_NOT_FOUND	RECORD_NOT_FOUND	✔	✔
//EVENT.CACHE_RETRIEVAL_TIMEOUT	CACHE_RETRIEVAL_TIMEOUT	✔
//EVENT.STORAGE_RETRIEVAL_TIMEOUT	STORAGE_RETRIEVAL_TIMEOUT	✔
//EVENT.CLOSED_SOCKET_INTERACTION	CLOSED_SOCKET_INTERACTION	✔
//EVENT.CLIENT_DISCONNECTED	CLIENT_DISCONNECTED	✔
//EVENT.INVALID_MESSAGE	INVALID_MESSAGE	✔
//EVENT.VERSION_EXISTS	VERSION_EXISTS	✔	✔
//EVENT.INVALID_VERSION	INVALID_VERSION	✔
//EVENT.PLUGIN_ERROR	PLUGIN_ERROR	✔
//EVENT.UNKNOWN_CALLEE	UNKNOWN_CALLEE	✔	✔

//Topic

//TopicConnection represents a connection related topic
const TopicConnection = "C"

//TopicAuth represents an auth related topic
const TopicAuth = "A"

//TopicError represents an error related topic
const TopicError = "X"

//TopicEvent represents an event related topic
const TopicEvent = "E"

//TopicRecord represents a record related topic
const TopicRecord = "R"

//TopicRPC represents an RPC related topic
const TopicRPC = "P"

//TopicPrivate represents a Private related topic
const TopicPrivate = "PRIVATE"

//Actions

const ActionAck = "A"
const ActionRead = "R"
const ActionRedirect = "RED"
const ActionChallenge = "CH"
const ActionChallengeResponse = "CHR"
const ActionCreate = "C"
const ActionUpdate = "U"
const ActionPatch = "P"
const ActionDelete = "D"
const ActionSubscribe = "S"
const ActionUnsubscribe = "uS"
const ActionHas = "H"
const ActionSnapshot = "SN"
const ActionListenSnapshot = "LSN"
const ActionListen = "L"
const ActionUnlisten = "UL"
const ActionListenAccept = "LA"
const ActionListenReject = "LR"
const ActionSubscriptionHasProvider = "SH"
const ActionSubscriptionsForPatternFound = "SF"
const ActionSubscriptionForPatternFound = "SP"
const ActionSubscriptionForPatternRemoved = "SR"
const ActionProviderUpdate = "PU"
const ActionQuery = "Q"
const ActionCreateOrRead = "CR"
const ActionEvent = "EVT"
const ActionError = "E"
const ActionRequest = "REQ"
const ActionResponse = "RES"
const ActionRejection = "REJ"
const ActionPing = "PI"
const ActionPong = "PO"

//Data Types

//DataType represents one of the available data types in an action
type DataType string

//TypesString indicates that the data in an action is of type string
const TypesString DataType = "S"

//TypesObject indicates that the data in an action is of type object (interface{})
const TypesObject DataType = "O"

//TypesNumber indicates that the data in an action is of type number
const TypesNumber DataType = "N"

//TypesNull indicates that the data in an action is nil
const TypesNull DataType = "L"

//TypesTrue indicates that the data in an action is true
const TypesTrue DataType = "T"

//TypesFalse indicates that the data in an action is false
const TypesFalse DataType = "F"

//TypesUndefined indicates that the data in an action is undefined
const TypesUndefined DataType = "U"
