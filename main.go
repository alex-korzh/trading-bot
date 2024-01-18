package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {

	stockSymbols, err := getSymbols()
	if err != nil {
		fmt.Println("Error during symbols retrieval:", err)
		os.Exit(1)
	}

	prices := make(map[string]float64)

	for _, stockSymbol := range stockSymbols {
		price, err := getPrice(stockSymbol)
		if err != nil {
			fmt.Println("Error during price retrieval:", err)
			os.Exit(1)
		}
		prices[stockSymbol] = price
	}

	fmt.Println("Last close prices of your stocks:")
	for symbol, price := range prices {
		fmt.Printf("%s: %.2f\n", symbol, price)
	}

}

func getSymbols() ([]string, error) {
	var symbols []string
	if len(os.Args) >= 2 {
		symbols = append(symbols, os.Args[1:]...)
	}

	// TODO add DB symbols to symbols slice

	return symbols, nil
}

type ClosePrice struct {
	Date  string  `json:"date"`
	Close float64 `json:"close"`
}

func getPrice(symbol string) (float64, error) {
	token := os.Getenv("TIINGO_API_TOKEN")
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
