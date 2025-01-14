package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = os.Getenv("POSTGRES_PORT")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

type Price struct {
	Id         int
	Name       string
	Price      float64
	CreateDate string
	Category   string
}

type Database struct {
	conn *sql.DB
}

func (db *Database) SetConnect() error {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	openedConn, err := sql.Open("postgres", conn)

	if err != nil {
		return err
	}

	db.conn = openedConn
	return nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

func (db *Database) InsertPrices(prices []Price) (int, int, error) {
	tx, err := db.conn.Begin()

	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
	}()

	query := "INSERT INTO prices (id, name, category, price, create_date) VALUES ($1, $2, $3, $4, $5);"

	for _, price := range prices {
		_, err = tx.Exec(query, price.Id, price.Name, price.Category, price.Price, price.CreateDate)

		if err != nil {
			return 0, 0, err
		}
	}

	sP, cC, err := getSumPriceAndCountCategories(tx)

	if tx.Commit() != nil {
		return 0, 0, err
	}

	return sP, cC, err
}

func (db *Database) GetPrices() ([]Price, error) {
	rows, err := db.conn.Query("SELECT id, name, category, price, create_date FROM prices;")

	if err != nil {
		return nil, err
	}

	var prices []Price
	var price Price

	for rows.Next() {
		err = rows.Scan(&price)

		if err != nil {
			continue
		}

		prices = append(prices, price)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return prices, nil
}

func getSumPriceAndCountCategories(tx *sql.Tx) (int, int, error) {
	row, err := tx.Query("SELECT SUM(price), COUNT(DISTINCT category) FROM prices;")

	if err != nil {
		return 0, 0, err
	}

	var sum, count int

	err = row.Scan(sum, count)

	if err != nil {
		return 0, 0, err
	}

	return sum, count, nil
}
