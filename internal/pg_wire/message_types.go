package pg_wire

type ClientMessageType byte

type ServerMessageType byte

type DescribeMessageType byte

const (
	ClientBind        ClientMessageType = 'B'
	ClientClose       ClientMessageType = 'C'
	ClientCopyData    ClientMessageType = 'd'
	ClientCopyDone    ClientMessageType = 'c'
	ClientCopyFail    ClientMessageType = 'f'
	ClientDescribe    ClientMessageType = 'D'
	ClientExecute     ClientMessageType = 'E'
	ClientFlush       ClientMessageType = 'H'
	ClientParse       ClientMessageType = 'P'
	ClientPassword    ClientMessageType = 'p'
	ClientSimpleQuery ClientMessageType = 'Q'
	ClientSync        ClientMessageType = 'S'
	ClientTerminate   ClientMessageType = 'X'

	ServerAuth                 ServerMessageType = 'R'
	ServerBindComplete         ServerMessageType = '2'
	ServerCommandComplete      ServerMessageType = 'C'
	ServerCloseComplete        ServerMessageType = '3'
	ServerCopyInResponse       ServerMessageType = 'G'
	ServerDataRow              ServerMessageType = 'D'
	ServerEmptyQuery           ServerMessageType = 'I'
	ServerErrorResponse        ServerMessageType = 'E'
	ServerNoticeResponse       ServerMessageType = 'N'
	ServerNoData               ServerMessageType = 'n'
	ServerParameterDescription ServerMessageType = 't'
	ServerParameterStatus      ServerMessageType = 'S'
	ServerParseComplete        ServerMessageType = '1'
	ServerPortalSuspended      ServerMessageType = 's'
	ServerReady                ServerMessageType = 'Z'
	ServerRowDescription       ServerMessageType = 'T'

	DescribePortal    DescribeMessageType = 'P'
	DescribeStatement DescribeMessageType = 'S'
)

var clientMessageName = map[ClientMessageType]string{
	ClientBind:        "Bind",
	ClientClose:       "Close",
	ClientCopyData:    "CopyData",
	ClientCopyDone:    "CopyDone",
	ClientCopyFail:    "CopyFail",
	ClientDescribe:    "Describe",
	ClientExecute:     "Execute",
	ClientFlush:       "Flush",
	ClientParse:       "Parse",
	ClientPassword:    "Password",
	ClientSimpleQuery: "SimpleQuery",
	ClientSync:        "Sync",
	ClientTerminate:   "Terminate",
}

var serverMessageName = map[ServerMessageType]string{
	ServerAuth:                 "Auth",
	ServerBindComplete:         "BindComplete",
	ServerCommandComplete:      "CommandComplete",
	ServerCloseComplete:        "CloseComplete",
	ServerCopyInResponse:       "CopyInResponse",
	ServerDataRow:              "DataRow",
	ServerEmptyQuery:           "EmptyQuery",
	ServerErrorResponse:        "ErrorResponse",
	ServerNoticeResponse:       "NoticeResponse",
	ServerNoData:               "NoData",
	ServerParameterDescription: "ParameterDescription",
	ServerParameterStatus:      "ParameterStatus",
	ServerParseComplete:        "ParseComplete",
	ServerPortalSuspended:      "PortalSuspended",
	ServerReady:                "Ready",
	ServerRowDescription:       "RowDescription",
}

var describeMessageType = map[DescribeMessageType]string{
	DescribePortal:    "Portal",
	DescribeStatement: "Statement",
}

func (c ClientMessageType) String() string {
	return clientMessageName[c]
}

func (s ServerMessageType) String() string {
	return serverMessageName[s]
}

func (d DescribeMessageType) String() string {
	return describeMessageType[d]
}
