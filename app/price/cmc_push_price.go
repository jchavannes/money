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
	Cancel      context.CancelFunc
	Conn        *websocket.Conn
	Investments []db.Investment
	ErrChan     chan error
}

func (c *CmcPushConn) Run() error {
	c.Ctx, c.Cancel = context.WithTimeout(context.Background(), 1*time.Minute)
	defer c.Cancel()

	var err error
	if c.Conn, _, err = websocket.Dial(c.Ctx, GetCmcPushPriceUrl(), nil); err != nil {
		return fmt.Errorf("error dialing websocket; %w", err)
	}
	defer c.Conn.CloseNow()

	go c.subscribe()
	go c.listen()

	select {
	case <-c.Ctx.Done():
		break
	case err := <-c.ErrChan:
		return fmt.Errorf("error processing websocket message; %w", err)
	}

	if err := c.Conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		return fmt.Errorf("error closing websocket normal; %w", err)
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
	var neededPrices = make([]uint, len(c.Investments))
	for i := range c.Investments {
		neededPrices[i] = c.Investments[i].Id
	}
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
		
		investmentPrice, err := GetInvestmentPriceFromCmcPushMessage(msg, c.Investments)
		if err != nil {
			c.ErrChan <- fmt.Errorf("error processing websocket message; %w", err)
			return
		}

		if investmentPrice == nil {
			continue
		}

		investmentPrice.Print()

		for i := 0; i < len(neededPrices); i++ {
			if neededPrices[i] == investmentPrice.InvestmentId {
				neededPrices = append(neededPrices[:i], neededPrices[i+1:]...)

				if err := investmentPrice.AddOrUpdate(); err != nil {
					c.ErrChan <- fmt.Errorf("error updating investment price: %#v; %w", investmentPrice, err)
					return
				}

				if len(neededPrices) == 0 {
					log.Println("All prices received, closing connection")
					c.Cancel()
					return
				}

				break
			}
		}
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
			break
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
