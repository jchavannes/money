package price_test

import (
	"encoding/json"
	"github.com/jchavannes/money/app/price"
	"log"
	"os"
	"testing"
)

func TestCmcPushPrice(t *testing.T) {
	cmcPushAck, err := os.ReadFile("./testdata/cmc_push_ack.json")
	if err != nil {
		t.Errorf("error reading cmc push ack; %v", err)
		return
	}

	var cmcPushAckJson price.CmcPushAckJson
	if err := json.Unmarshal(cmcPushAck, &cmcPushAckJson); err != nil {
		t.Errorf("error parsing cmc push ack; %v", err)
		return
	}

	log.Printf("cmc push ack: %#v\n", cmcPushAckJson)
}