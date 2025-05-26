package pg

import "encoding/binary"

type Message struct {
	Type   byte
	Length int32
	Data   []byte
}

func (m *Message) Binary() []byte {
	b := make([]byte, m.Length+1)
	b[0] = m.Type
	binary.BigEndian.PutUint32(b[1:], uint32(m.Length))
	copy(b[5:], m.Data)
	return b
}

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

func (m ClientMessageType) ToString() string {
	switch m {
	case ClientBind:
		return "Bind"
	case ClientClose:
		return "Close"
	case ClientCopyData:
		return "CopyData"
	case ClientCopyDone:
		return "CopyDone"
	case ClientCopyFail:
		return "CopyFail"
	case ClientDescribe:
		return "Describe"
	case ClientExecute:
		return "Execute"
	case ClientFlush:
		return "Flush"
	case ClientParse:
		return "Parse"
	case ClientPassword:
		return "Password"
	case ClientSimpleQuery:
		return "SimpleQuery"
	case ClientSync:
		return "Sync"
	case ClientTerminate:
		return "Terminate"
	default:
		return "Unknown"
	}
}

func (m ServerMessageType) ToString() string {
	switch m {
	case ServerAuth:
		return "Auth"
	case ServerBindComplete:
		return "BindComplete"
	case ServerCommandComplete:
		return "CommandComplete"
	case ServerCloseComplete:
		return "CloseComplete"
	case ServerCopyInResponse:
		return "CopyInResponse"
	case ServerDataRow:
		return "DataRow"
	case ServerEmptyQuery:
		return "EmptyQuery"
	case ServerErrorResponse:
		return "ErrorResponse"
	case ServerNoticeResponse:
		return "NoticeResponse"
	case ServerNoData:
		return "NoData"
	case ServerParameterDescription:
		return "ParameterDescription"
	case ServerParameterStatus:
		return "ParameterStatus"
	case ServerParseComplete:
		return "ParseComplete"
	case ServerPortalSuspended:
		return "PortalSuspended"
	case ServerReady:
		return "Ready"
	case ServerRowDescription:
		return "RowDescription"
	default:
		return "Unknown"
	}
}

func (m DescribeMessageType) ToString() string {
	switch m {
	case DescribePortal:
		return "Portal"
	case DescribeStatement:
		return "Statement"
	default:
		return "Unknown"
	}
}
