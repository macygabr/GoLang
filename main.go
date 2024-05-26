package main

import (
	"golang/cash"
	"golang/database"
	"golang/server"
)

func main() {
	db := new(database.DataBase)
	defer db.Connect().Unsubscribe()

	cash := new(cash.Cash)
	defer cash.Regenerate().Unsubscribe()

	server := new(server.Server)
	defer server.Start().Unsubscribe()
}
