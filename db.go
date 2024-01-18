package main

import (
	"database/sql"
	"fmt"
	"math"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "data.db"
const initSql = `
CREATE TABLE IF NOT EXISTS prices (
	symbol TEXT NOT NULL PRIMARY KEY,
	date DATETIME NOT NULL,
	price REAL NOT NULL
);`

type Repository struct {
	db *sql.DB
}

func newRepository() (*Repository, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	r := Repository{db: db}
	return &r, nil
}

func (r *Repository) initDB() error {
	if _, err := r.db.Exec(initSql); err != nil {
		return err
	}
	return nil
}

func (r *Repository) getDBSymbols() ([]string, error) {

	rows, err := r.db.Query("SELECT symbol FROM prices")
	if err != nil {
		return nil, err
	}
	var res []string
	var symbol string
	for rows.Next() {
		err := rows.Scan(&symbol)
		if err != nil {
			return nil, err
		}
		res = append(res, symbol)
	}
	return res, nil
}

func (r *Repository) insertSymbols(symbols map[string]float64) error {
	rows, err := r.db.Query("SELECT symbol, price FROM prices")
	if err != nil {
		return err
	}
	var symbol string
	var price float64
	for rows.Next() {
		err := rows.Scan(&symbol, &price)
		if err != nil {
			return err
		}
		if v, ok := symbols[symbol]; ok {
			if math.Abs(v-price) > 0.01 {
				_, err := r.db.Exec("UPDATE prices SET price=?, date=CURRENT_TIMESTAMP WHERE symbol=?", v, symbol)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("Skipped updating %s\n", symbol)
			}
			delete(symbols, symbol)
		}
	}
	for s, p := range symbols {
		_, err := r.db.Exec("INSERT INTO prices VALUES(?, CURRENT_TIMESTAMP, ?)", s, p)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *Repository) getPrices() (map[string]float64, error) {
	rows, err := r.db.Query("SELECT symbol, price FROM prices")
	if err != nil {
		return nil, err
	}
	res := make(map[string]float64)
	var symbol string
	var price float64
	for rows.Next() {
		err := rows.Scan(&symbol, &price)
		if err != nil {
			return nil, err
		}
		res[symbol] = price
	}
	return res, nil
}
