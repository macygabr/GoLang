package main

import (
	"golang/cash"
	"golang/database"
	"golang/server"
)

func main() {
	db := database.NewDataBase()
	defer db.Connect().Unsubscribe()

	cash := cash.NewCash()
	defer cash.Regenerate().Unsubscribe()

	server := server.NewServer(db, nil)
	defer server.Start().Unsubscribe()
}
