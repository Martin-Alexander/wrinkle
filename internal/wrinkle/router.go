package wrinkle

import (
	"context"
	"io"
	"net"
	"time"
	"wrinkle/internal/pg_wire"
)

type Router struct {
	reader     RouterMessageReader
	writer     RouterMessageWriter
	controller *Controller
	errorCh    chan<- error
	ctx        context.Context
	cancel     context.CancelFunc
}

type RouterMessageReader interface {
	ReadMessage(conn io.Reader, sender pg_wire.Sender) (*pg_wire.Message, error)
}

type RouterMessageWriter interface {
	WriteMessage(conn io.Writer, message *pg_wire.Message) error
}

func NewRouter(
	reader RouterMessageReader,
	writer RouterMessageWriter,
	controller *Controller,
	errorCh chan<- error,
) *Router {
	ctx, cancel := context.WithCancel(context.Background())

	return &Router{
		reader:     reader,
		writer:     writer,
		controller: controller,
		errorCh:    errorCh,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (r *Router) Start(feConn net.Conn, beConn net.Conn) {
	go r.readMessageLoop(feConn, r.controller.inboundClientMsgCh, pg_wire.Frontend)
	go r.readMessageLoop(beConn, r.controller.inboundDbMsgCh, pg_wire.Backend)
	go r.writeMessageLoop(feConn, r.controller.outboundClientMsgCh)
	go r.writeMessageLoop(beConn, r.controller.outboundDbMsgCh)
}

func (r *Router) Stop() {
	r.cancel()
}

func (r *Router) readMessageLoop(conn net.Conn, msgRecvCh chan<- *pg_wire.Message, sender pg_wire.Sender) {
	for {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		message, err := r.reader.ReadMessage(conn, sender)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && !netErr.Timeout() {
				r.errorCh <- err
			}

			continue
		}

		conn.SetReadDeadline(time.Time{})

		select {
		case msgRecvCh <- message:
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Router) writeMessageLoop(conn net.Conn, sendMsgCh <-chan *pg_wire.Message) {
	for {
		select {
		case message := <-sendMsgCh:
			if err := r.writer.WriteMessage(conn, message); err != nil {
				r.errorCh <- err
				continue
			}

		case <-r.ctx.Done():
			return
		}
	}
}
