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

type Bridge struct {
	frontendConn net.Conn
	backendConn  net.Conn
	ctx          context.Context
	cancel       context.CancelCauseFunc
}

func NewBridge(
	frontendConn net.Conn,
	backendConn net.Conn,
) *Bridge {
	return &Bridge{
		frontendConn: frontendConn,
		backendConn:  backendConn,
	}
}

func (b *Bridge) StartRelaying() error {
	b.ctx, b.cancel = context.WithCancelCause(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	go b.startForwardRelaying(b.frontendConn, b.backendConn, &wg)
	go b.startForwardRelaying(b.backendConn, b.frontendConn, &wg)

	wg.Wait()

	return context.Cause(b.ctx)
}

func (b *Bridge) startForwardRelaying(source net.Conn, destination net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
		}

		buffer := make([]byte, 4096)

		n, err := source.Read(buffer)
		if errors.Is(err, io.EOF) && n == 0 {
			return
		}
		if err != nil {
			b.cancel(err)
			return
		}

		logPacket("Packet", buffer[:n])

		if _, err := destination.Write(buffer[:n]); err != nil {
			b.cancel(err)
			return
		}
	}
}

func logPacket(label string, data []byte) {
	log.Printf(" -- [%s]:\n%s\n", label, hex.Dump(data))
}
