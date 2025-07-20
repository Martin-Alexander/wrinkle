package proxy

import (
	"context"
	"encoding/hex"
	"io"
	"log"
	"net"
	"time"
)

type Broker struct {
	reader      BrokerMessageReader
	writer      BrokerMessageWriter
	feMsgRecvCh chan<- []byte
	beMsgRecvCh chan<- []byte
	sendFeMsgCh <-chan []byte
	sendBeMsgCh <-chan []byte
	errorCh     chan<- error
	ctx         context.Context
	cancel      context.CancelFunc
}

type BrokerMessageReader interface {
	ReadMessage(conn io.Reader) ([]byte, error)
}

type BrokerMessageWriter interface {
	WriteMessage(conn io.Writer, message []byte) error
}

func NewBroker(
	reader BrokerMessageReader,
	writer BrokerMessageWriter,
	feMsgRecvCh chan<- []byte,
	beMsgRecvCh chan<- []byte,
	sendFeMsgCh <-chan []byte,
	sendBeMsgCh <-chan []byte,
	errorCh chan<- error,
) *Broker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Broker{
		reader:      reader,
		writer:      writer,
		feMsgRecvCh: feMsgRecvCh,
		beMsgRecvCh: beMsgRecvCh,
		sendFeMsgCh: sendFeMsgCh,
		sendBeMsgCh: sendBeMsgCh,
		errorCh:     errorCh,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (b *Broker) Start(feConn net.Conn, beConn net.Conn) {
	go b.readLoop(feConn, b.feMsgRecvCh)
	go b.readLoop(beConn, b.beMsgRecvCh)
	go b.writeLoop(feConn, b.sendFeMsgCh)
	go b.writeLoop(beConn, b.sendBeMsgCh)
}

func (b *Broker) Stop() {
	b.cancel()
}

func (b *Broker) readLoop(conn net.Conn, msgRecvCh chan<- []byte) {
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
			log.Printf("\n%s\n", hex.Dump(message))
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *Broker) writeLoop(conn net.Conn, sendMsgCh <-chan []byte) {
	for {
		select {
		case message := <-sendMsgCh:
			if err := b.writer.WriteMessage(conn, message); err != nil {
				b.errorCh <- err
				continue
			}

			log.Printf("\n%s\n", hex.Dump(message))

		case <-b.ctx.Done():
			return
		}
	}
}
