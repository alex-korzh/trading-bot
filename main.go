package main

import (
	"fmt"
	"os"
)

func main() {
	repo, err := newRepository()
	if err != nil {
		fmt.Println("Error during DB connection:", err)
		os.Exit(1)
	}

	if err := repo.initDB(); err != nil {
		fmt.Println("Error during DB initialization:", err)
		os.Exit(1)
	}

	stockSymbols, err := getSymbols(repo)
	if err != nil {
		fmt.Println("Error during symbols retrieval:", err)
		os.Exit(1)
	}
	if len(stockSymbols) == 0 {
		fmt.Println("No stocks in DB, and no stocks provided. Provide stocks symbols.")
		os.Exit(1)
	}

	prices := make(map[string]float64)

	for stockSymbol := range stockSymbols {
		price, err := getPrice(stockSymbol)
		if err != nil {
			fmt.Println("Error during price retrieval:", err)
			os.Exit(1)
		}
		prices[stockSymbol] = price
	}

	err = repo.insertSymbols(prices)
	if err != nil {
		fmt.Println("Error during price saving attempt:", err)
		os.Exit(1)
	}

	prices, err = repo.getPrices()
	if err != nil {
		fmt.Println("Error during price retrieval from DB:", err)
		os.Exit(1)
	}

	fmt.Println("Last close prices of your stocks:")
	for symbol, price := range prices {
		fmt.Printf("%s: %.2f\n", symbol, price)
	}

	defer repo.db.Close()

}

func getSymbols(repo *Repository) (map[string]bool, error) {
	var symbols []string
	if len(os.Args) >= 2 {
		symbols = append(symbols, os.Args[1:]...)
	}

	dbSymbols, err := repo.getDBSymbols()
	if err != nil {
		return nil, err
	}

	symbols = append(symbols, dbSymbols...)

	symbolsMap := make(map[string]bool)

	for _, s := range symbols {
		symbolsMap[s] = true
	}

	return symbolsMap, nil
}
