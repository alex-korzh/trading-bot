package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ClosePrice struct {
	Date  string  `json:"date"`
	Close float64 `json:"close"`
}

func getPrice(symbol string) (float64, error) {
	token := os.Getenv("TIINGO_API_TOKEN")
	if token == "" {
		return 0, fmt.Errorf("TIINGO_API_TOKEN env variable is not set")
	}
	uri := fmt.Sprintf("https://api.tiingo.com/tiingo/daily/%s/prices?token=%s&columns=close", symbol, token)
	resp, err := http.Get(uri)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var priceObj []ClosePrice
	err = json.NewDecoder(resp.Body).Decode(&priceObj)
	if err != nil {
		return 0, err
	}
	price := priceObj[0].Close
	return price, nil
}
