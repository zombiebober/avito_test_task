package main

import (
	"log"
	"os"
)

func main() {

	a := App{}
	log.Print("Сервер запущен")
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
	os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	a.Run(":8080")
}
