package wrinkle

import (
	"context"
	"encoding/hex"
	"io"
	"log"
	"net"
	"time"
	"wrinkle/internal/pg_wire"
)

type Broker struct {
	reader     BrokerMessageReader
	writer     BrokerMessageWriter
	controller *Controller
	errorCh    chan<- error
	ctx        context.Context
	cancel     context.CancelFunc
}

type BrokerMessageReader interface {
	ReadMessage(conn io.Reader) (*pg_wire.Message, error)
}

type BrokerMessageWriter interface {
	WriteMessage(conn io.Writer, message *pg_wire.Message) error
}

func NewBroker(
	reader BrokerMessageReader,
	writer BrokerMessageWriter,
	controller *Controller,
	errorCh chan<- error,
) *Broker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Broker{
		reader:     reader,
		writer:     writer,
		controller: controller,
		errorCh:    errorCh,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (b *Broker) Start(feConn net.Conn, beConn net.Conn) {
	go b.readLoop(feConn, b.controller.inboundClientMsgCh)
	go b.readLoop(beConn, b.controller.inboundDbMsgCh)
	go b.writeLoop(feConn, b.controller.outboundClientMsgCh)
	go b.writeLoop(beConn, b.controller.outboundDbMsgCh)
}

func (b *Broker) Stop() {
	b.cancel()
}

func (b *Broker) readLoop(conn net.Conn, msgRecvCh chan<- *pg_wire.Message) {
	for {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		message, err := b.reader.ReadMessage(conn)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && !netErr.Timeout() {
				b.errorCh <- err
			}

			continue
		}

		conn.SetReadDeadline(time.Time{})

		select {
		case msgRecvCh <- message:
			log.Printf("\n%s\n", hex.Dump(message.Bytes()))
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *Broker) writeLoop(conn net.Conn, sendMsgCh <-chan *pg_wire.Message) {
	for {
		select {
		case message := <-sendMsgCh:
			if err := b.writer.WriteMessage(conn, message); err != nil {
				b.errorCh <- err
				continue
			}

			log.Printf("\n%s\n", hex.Dump(message.Bytes()))

		case <-b.ctx.Done():
			return
		}
	}
}
