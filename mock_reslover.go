package ipm

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func genRsp() string {
	rand.Seed(time.Now().UnixNano())
	min := 39000
	max := 41000
	amount := rand.Intn(max-min) + min
	return fmt.Sprintf(`{ "amount": %d }`, amount)
}

type MockPriceResolver struct {
}

func (MockPriceResolver) Get(target string) (time.Time, Decimal, error) {
	resp := genRsp()

	data := struct {
		Amount Decimal
	}{}

	err := json.Unmarshal([]byte(resp), &data)

	return time.Now(), data.Amount, err
}
