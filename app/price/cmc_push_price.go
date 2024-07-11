package price

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jchavannes/money/app/db"
	"log"
	"nhooyr.io/websocket"
	"strconv"
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
	Price float32 `json:"p"`
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
	// TODO: Parse different message types and save prices.
	// TODO: Once each investment has had a price update, close.
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

func GetInvestmentPriceFromCmcPushMessage(msg []byte, investments []db.Investment) (*db.InvestmentPrice, error) {
	var cmcPushPriceJson CmcPushPriceJson
	if err := json.Unmarshal(msg, &cmcPushPriceJson); err != nil {
		return nil, fmt.Errorf("error parsing cmc push price; %w", err)
	}

	if cmcPushPriceJson.D.Id == 0 {
		// Different socket message
		return nil, nil
	}

	timeInt, err := strconv.Atoi(cmcPushPriceJson.Time)
	if err != nil {
		return nil, fmt.Errorf("error converting time to int; %w", err)
	}

	var foundInvestment *db.Investment
	for _, investment := range investments {
		if cmcPushPriceJson.D.Id == GetIdFromSymbol(investment.Symbol) {
			foundInvestment = &investment
		}
	}
	if foundInvestment == nil {
		return nil, fmt.Errorf("unable to find investment for price id: %d", cmcPushPriceJson.D.Id)
	}

	return &db.InvestmentPrice{
		Investment:   *foundInvestment,
		InvestmentId: foundInvestment.Id,
		Timestamp:    time.UnixMilli(int64(timeInt)).Unix(),
		Price:        cmcPushPriceJson.D.Price,
	}, nil
}
