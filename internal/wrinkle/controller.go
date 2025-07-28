package wrinkle

import (
	"context"
	"encoding/hex"
	"log"
	"wrinkle/internal/pg_wire"
)

type Controller struct {
	inboundClientMsgCh  chan *pg_wire.Message
	inboundDbMsgCh      chan *pg_wire.Message
	inboundBrainMsgCh   chan *pg_wire.Message
	outboundClientMsgCh chan *pg_wire.Message
	outboundDbMsgCh     chan *pg_wire.Message
	outboundBrainMsgCh  chan *pg_wire.Message
	ctx                 context.Context
	cancel              context.CancelFunc
}

func NewController() *Controller {
	ctx, cancel := context.WithCancel(context.Background())

	return &Controller{
		inboundClientMsgCh:  make(chan *pg_wire.Message, 100),
		inboundDbMsgCh:      make(chan *pg_wire.Message, 100),
		inboundBrainMsgCh:   make(chan *pg_wire.Message, 100),
		outboundClientMsgCh: make(chan *pg_wire.Message, 100),
		outboundDbMsgCh:     make(chan *pg_wire.Message, 100),
		outboundBrainMsgCh:  make(chan *pg_wire.Message, 100),
		ctx:                 ctx,
		cancel:              cancel,
	}
}

func (c *Controller) InboundClientMsgCh() chan<- *pg_wire.Message {
	return c.inboundClientMsgCh
}

func (c *Controller) InboundDbMsgCh() chan<- *pg_wire.Message {
	return c.inboundDbMsgCh
}

func (c *Controller) InboundBrainMsgCh() chan<- *pg_wire.Message {
	return c.inboundBrainMsgCh
}

func (c *Controller) OutboundClientMsgCh() <-chan *pg_wire.Message {
	return c.outboundClientMsgCh
}

func (c *Controller) OutboundDbMsgCh() <-chan *pg_wire.Message {
	return c.outboundDbMsgCh
}

func (c *Controller) OutboundBrainMsgCh() <-chan *pg_wire.Message {
	return c.outboundBrainMsgCh
}

func (c *Controller) Start() {
	go c.readLoop(c.inboundClientMsgCh, c.HandleClientMessage)
	go c.readLoop(c.inboundDbMsgCh, c.HandleDbMessage)
	go c.readLoop(c.inboundBrainMsgCh, c.HandleBrainMessage)
}

func (c *Controller) Stop() {
	c.cancel()
}

func (c *Controller) readLoop(channel <-chan *pg_wire.Message, handler func(msg *pg_wire.Message)) {
	for {
		select {
		case msg := <-channel:
			handler(msg)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Controller) HandleClientMessage(msg *pg_wire.Message) {
	log.Printf("[Client] %s\n%s\n", msg.Name(), hex.Dump(msg.Bytes()))

	c.outboundDbMsgCh <- msg // For now, just send everything to the DB
}

func (c *Controller) HandleDbMessage(msg *pg_wire.Message) {
	log.Printf("[Server] %s\n%s\n", msg.Name(), hex.Dump(msg.Bytes()))

	c.outboundClientMsgCh <- msg
}

func (c *Controller) HandleBrainMessage(msg *pg_wire.Message) {
	c.outboundClientMsgCh <- msg
}
