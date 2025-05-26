package proxy

import (
	"context"
	"encoding/hex"
	"log"
	"net"
	"sync"
	"wrinkle/internal/pg"
)

type Relay struct {
	frontendConn net.Conn
	backendConn  net.Conn
	ctx          context.Context
	cancel       context.CancelCauseFunc
}

func NewRelay(
	frontendConn net.Conn,
	backendConn net.Conn,
) *Relay {
	return &Relay{
		frontendConn: frontendConn,
		backendConn:  backendConn,
	}
}

func (r *Relay) Start() error {
	r.ctx, r.cancel = context.WithCancelCause(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	go r.startForwardingFrontend(r.frontendConn, r.backendConn, &wg)
	go r.startForwardingBackend(r.backendConn, r.frontendConn, &wg)

	wg.Wait()

	return context.Cause(r.ctx)
}

func (r *Relay) startForwardingFrontend(source net.Conn, destination net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}

		message, err := ReadMessage(source)
		if err != nil {
			r.cancel(err)
			return
		}

		log.Printf(
			"-- %s -- %s:\n%s\n",
			"Client",
			pg.ClientMessageType(message.Type).ToString(),
			hex.Dump(message.Binary()),
		)

		if _, err := destination.Write(message.Binary()); err != nil {
			r.cancel(err)
			return
		}
	}
}

func (r *Relay) startForwardingBackend(source net.Conn, destination net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}

		message, err := ReadMessage(source)
		if err != nil {
			r.cancel(err)
			return
		}

		log.Printf(
			"-- %s -- %s:\n%s\n",
			"Server",
			pg.ServerMessageType(message.Type).ToString(),
			hex.Dump(message.Binary()),
		)

		if _, err := destination.Write(message.Binary()); err != nil {
			r.cancel(err)
			return
		}
	}
}
