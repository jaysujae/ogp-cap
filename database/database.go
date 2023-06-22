package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5455
	user   = "postgresUser"
	dbname = "postgresDB"
	password = "postgresPW"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	_, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Postgres Database has been setup successfully")
}
