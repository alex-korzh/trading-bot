package main

import (
	"math/rand"
	"time"
)

type Decision string

const (
	doNothing Decision = "DO_NOTHING"
	buy       Decision = "BUY"
	sell      Decision = "SELL"
)

type PricePoint struct {
	price float64
	date  time.Time
}

func RandomStrategy(symbol string, prices ...PricePoint) (Decision, error) {
	choices := [3]Decision{doNothing, buy, sell}
	randomChoice := rand.Intn(3)
	return choices[randomChoice], nil
}
