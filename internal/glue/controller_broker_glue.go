package glue

import (
	"log/slog"
	"wrinkle/internal/pg_wire"
	"wrinkle/internal/wrinkle"
)

type ControllerBrokerGlue struct {
	FeMsgRecvCh chan<- []byte
	BeMsgRecvCh chan<- []byte
	SendFeMsgCh <-chan []byte
	SendBeMsgCh <-chan []byte
}

func NewControllerBrokerGlue(controller *wrinkle.Controller) *ControllerBrokerGlue {
	feMsgRecvCh := make(chan []byte, 100)
	beMsgRecvCh := make(chan []byte, 100)
	sendFeMsgCh := make(chan []byte, 100)
	sendBeMsgCh := make(chan []byte, 100)

	glue := &ControllerBrokerGlue{
		FeMsgRecvCh: feMsgRecvCh,
		BeMsgRecvCh: beMsgRecvCh,
		SendFeMsgCh: sendFeMsgCh,
		SendBeMsgCh: sendBeMsgCh,
	}

	go glue.parseDataLoop(feMsgRecvCh, controller.InboundClientMsgCh())
	go glue.parseDataLoop(beMsgRecvCh, controller.InboundDbMsgCh())
	go glue.serializeMessageLoop(controller.OutboundClientMsgCh(), sendFeMsgCh)
	go glue.serializeMessageLoop(controller.OutboundDbMsgCh(), sendBeMsgCh)

	return glue
}

func (g *ControllerBrokerGlue) parseDataLoop(
	byteCh <-chan []byte,
	msgCh chan<- *pg_wire.Message,
) {
	for data := range byteCh {
		msg, err := pg_wire.FromBytes(data)
		if err != nil {
			slog.Error("Failed to parse message", "data", data, "error", err)
			continue
		}

		msgCh <- msg
	}
}

func (g *ControllerBrokerGlue) serializeMessageLoop(
	msgCh <-chan *pg_wire.Message,
	byteCh chan<- []byte,
) {
	for msg := range msgCh {
		byteCh <- msg.Bytes()
	}
}
