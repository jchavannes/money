package price

import (
	"context"
	"fmt"
	"github.com/jchavannes/money/app/db"
	"log"
	"nhooyr.io/websocket"
	"strconv"
	"strings"
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

func CmcPushPrice(investments []db.Investment) error {
	url := "wss://push.coinmarketcap.com/ws?device=web&client_source=coin_detail_page"
	var idStrings = make([]string, len(investments))
	for i := range investments {
		idStrings[i] = strconv.Itoa(GetIdFromSymbol(investments[i].Symbol))
	}
	log.Printf("id strings: %s\n", strings.Join(idStrings, ","))
	sendMessage := `{"method":"RSUBSCRIPTION","params":["main-site@crypto_price_5s@{}@normal","` + strings.Join(idStrings, ",") + `"]}`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	log.Println("Connecting to websocket...")
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("error dialing websocket; %w", err)
	}
	defer conn.CloseNow()

	var errChan = make(chan error)

	go func() {
		time.Sleep(1 * time.Second)
		log.Println("Sending websocket subscribe message")
		if err := conn.Write(ctx, websocket.MessageText, []byte(sendMessage)); err != nil {
			errChan <- fmt.Errorf("error writing websocket message; %w", err)
		}
	}()

	go func() {
		for {
			msgType, msg, err := conn.Read(ctx)
			if err != nil {
				errChan <- fmt.Errorf("error reading websocket message; %w", err)
				return
			}
			if msgType != websocket.MessageText {
				errChan <- fmt.Errorf("unexpected message type; %w", err)
				return
			}
			log.Printf("Received message: %s\n", string(msg))
		}
	}()

	select {
	case <-ctx.Done():
		break
	case err := <-errChan:
		return fmt.Errorf("error processing websocket message; %w", err)
	}

	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		return fmt.Errorf("error closing websocket; %w", err)
	}

	return nil
}
