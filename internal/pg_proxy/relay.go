package pg_proxy

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net"
	"sync"
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

	go r.startForwarding(r.frontendConn, r.backendConn, &wg)
	go r.startForwarding(r.backendConn, r.frontendConn, &wg)

	wg.Wait()

	return context.Cause(r.ctx)
}

func (r *Relay) startForwarding(source net.Conn, destination net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}

		buffer := make([]byte, 4096)

		n, err := source.Read(buffer)
		if errors.Is(err, io.EOF) && n == 0 {
			return
		}
		if err != nil {
			r.cancel(err)
			return
		}

		logPacket("Packet", buffer[:n])

		if _, err := destination.Write(buffer[:n]); err != nil {
			r.cancel(err)
			return
		}
	}
}

func logPacket(label string, data []byte) {
	log.Printf(" -- [%s]:\n%s\n", label, hex.Dump(data))
}
