package db

import (
	"bytes"
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

func connect() (*sql.DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	return sql.Open("postgres", conn)
}

func InsertPrices(prices []Price) error {
	db, err := connect()

	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	var buffer bytes.Buffer

	buffer.WriteString("INSERT INTO prices (id, name, category, price, create_date) VALUES ")

	for idx, price := range prices {
		buffer.WriteString(
			fmt.Sprintf(
				"('%d', '%s', '%s', '%f', '%s')",
				price.Id,
				price.Name,
				price.Category,
				price.Price,
				price.CreateDate,
			),
		)
		if idx != len(prices)-1 {
			buffer.WriteString(",")
		} else {
			buffer.WriteString(";")
		}
	}

	_, err = db.Exec(buffer.String())

	return err
}

// Запрос на получение (БД)
func GetPrices() ([]Price, error) {
	db, err := connect()

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	rows, err := db.Query("SELECT id, name, category, price, create_date FROM prices;")

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

	return prices, nil
}

func GetSumPriceAndCountCategories() (int, int, error) {
	db, err := connect()

	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	row, err := db.Query("SELECT SUM(price), COUNT(DISTINCT category) FROM prices;")

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
