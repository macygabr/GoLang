package main

import (
	"fmt"
	"golang/cash"
	"golang/database"
	"golang/server"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	db := database.NewDataBase()
	db.Connect()

	cash := cash.NewCash()
	cash.Regenerate()

	//_______________________________________________________________________________________________________________
	sc, _ := stan.Connect("test-cluster", "client_id", stan.NatsURL("nats://0.0.0.0:4222"))
	sub, _ := sc.Subscribe("parseFile", func(msg *stan.Msg) {
		fmt.Printf("Received message!: %s\n", string(msg.Data))
		db.ReadFile()
	})
	defer sub.Unsubscribe()
	time.Sleep(time.Second)
	//_______________________________________________________________________________________________________________

	server := server.NewServer(db, nil)
	server.Start()
}

func CreateSubscribe() {
	// //_______________________________________________________________________________________________________________
	// sc, _ := stan.Connect("test-cluster", "client_id", stan.NatsURL("nats://0.0.0.0:4222"))
	// sub, _ := sc.Subscribe("parseFile", func(msg *stan.Msg) {
	// 	fmt.Printf("Received message!: %s\n", string(msg.Data))
	// 	db.ReadFile()
	// })
	// defer sub.Unsubscribe()
	// time.Sleep(time.Second)
	// //_______________________________________________________________________________________________________________
}
