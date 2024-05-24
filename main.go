package main

import (
	"level0/database"
	"level0/server"
)

func main() {
	db := new(database.DataBase)
	db.Connect()

	server := new(server.Server)
	server.Start(db)
}
