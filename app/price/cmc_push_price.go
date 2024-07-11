package price

import (
	"context"
	"fmt"
	"github.com/jchavannes/money/app/db"
	"log"
	"nhooyr.io/websocket"
	"time"
)

type CmcPushAckJson struct {
	Id   int    `json:"id"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CmcPushPriceJson struct {
	D    CmcPushPriceJsonD `json:"d"`
	Time string            `json:"t"` // eg. "1720661812642"
}

type CmcPushPriceJsonD struct {
	Id    int     `json:"id"`
	Price float64 `json:"p"`
}

func RunCmcPushPrice(investments []db.Investment) error {
	var pc = &CmcPushConn{
		Investments: investments,
		ErrChan:     make(chan error),
	}

	if err := pc.Run(); err != nil {
		return fmt.Errorf("error subscribing and listening to websocket; %w", err)
	}
	return nil
}

type CmcPushConn struct {
	Ctx         context.Context
	Conn        *websocket.Conn
	Investments []db.Investment
	ErrChan     chan error
}

func (c *CmcPushConn) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	c.Ctx = ctx

	var err error
	if c.Conn, _, err = websocket.Dial(ctx, GetCmcPushPriceUrl(), nil); err != nil {
		return fmt.Errorf("error dialing websocket; %w", err)
	}
	defer func() {
		if err := c.Conn.CloseNow(); err != nil {
			log.Printf("error closing websocket; %v", err)
		}
	}()

	go c.subscribe()
	go c.listen()

	select {
	case <-ctx.Done():
		break
	case err := <-c.ErrChan:
		return fmt.Errorf("error processing websocket message; %w", err)
	}

	if err := c.Conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		return fmt.Errorf("error closing websocket; %w", err)
	}

	return nil
}

func (c *CmcPushConn) subscribe() {
	log.Println("Sending websocket subscribe message")
	msg := GetCmcPushPriceSubscribeMessage(c.Investments)
	if err := c.Conn.Write(c.Ctx, websocket.MessageText, []byte(msg)); err != nil {
		c.ErrChan <- fmt.Errorf("error writing websocket message; %w", err)
	}
}

func (c *CmcPushConn) listen() {
	for {
		msgType, msg, err := c.Conn.Read(c.Ctx)
		if err != nil {
			c.ErrChan <- fmt.Errorf("error reading websocket message; %w", err)
			return
		}
		if msgType != websocket.MessageText {
			c.ErrChan <- fmt.Errorf("unexpected message type; %w", err)
			return
		}
		log.Printf("Received message: %s\n", string(msg))
	}
}
