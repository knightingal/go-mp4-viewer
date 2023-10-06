package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func initDB() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "000000",
		Addr:                 "localhost:3306",
		DBName:               "mp4viewer",
		AllowNativePasswords: true,
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")
}
