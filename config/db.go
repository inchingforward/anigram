package config

import (
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB(datasource string) {
	x, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		log.Fatal(err)
	}

	x.SetMaxOpenConns(2)
	x.SetMaxIdleConns(2)

	DB = x

	fmt.Println("Connected to anigram db")
}
