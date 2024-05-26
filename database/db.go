package database

import (
	"database/sql"
	"encoding/json"
	"golang/model/task"
	"golang/model/user"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
)

type DataBase struct {
	db *sql.DB
	sc stan.Conn
}

func NewDataBase() *DataBase {
	return &DataBase{nil, nil}
}

func (v *DataBase) Connect() stan.Subscription {
	var sub = v.Listen()
	conninfo := "user=postgres password=postgres host=127.0.0.1 sslmode=disable"
	db, err := sql.Open("postgres", conninfo)

	if err != nil {
		log.Fatal(err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS delivery ( Name varchar(255), Phone varchar(255), Zip varchar(255), City varchar(255), Address varchar(255), Region varchar(255), Email varchar(255) )")
	db.Exec("CREATE TABLE IF NOT EXISTS payment ( Transaction varchar(255) , RequestID varchar(255), Currency varchar(255), Provider varchar(255), Amount INTEGER, PaymentDt INTEGER, Bank varchar(255), DeliveryCost INTEGER, GoodsTotal INTEGER, CustomFee INTEGER)")
	db.Exec("CREATE TABLE IF NOT EXISTS items ( ChrtID INTEGER , TrackNumber varchar(255), Price INTEGER, Rid varchar(255), Name varchar(255), Sale INTEGER, Size varchar(255), TotalPrice INTEGER, NmID INTEGER, Brand varchar(255), Status INTEGER)")
	db.Exec("CREATE TABLE IF NOT EXISTS orders ( OrderUID varchar(255), TrackNumber varchar(255), Entry varchar(255), Locale varchar(255), InternalSignature varchar(255), CustomerID varchar(255), DeliveryService varchar(255), Shardkey varchar(255), SmID INTEGER, DateCreated varchar(255), OofShard varchar(255))")

	v.db = db
	return sub
}

func (v *DataBase) ReadFile(name string) {
	jsonData, err := ioutil.ReadFile("server/download/" + name)
	if err != nil {
		log.Fatal(err)
	}

	if !json.Valid(jsonData) {
		return
	}

	var data user.UserData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal(err)
	}
	v.insertDelivery(data)
	v.insertPayment(data)
	v.insertItems(data)
	v.insertOrders(data)

}

func (v *DataBase) insertDelivery(data user.UserData) {
	stmt, err := v.db.Prepare("INSERT INTO delivery (Name, Phone, Zip, City, Address, Region, Email) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.Delivery.Name, data.Delivery.Phone, data.Delivery.Zip, data.Delivery.City, data.Delivery.Address, data.Delivery.Region, data.Delivery.Email)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *DataBase) insertPayment(data user.UserData) {
	stmt, err := v.db.Prepare("INSERT INTO payment (Transaction, RequestID, Currency, Provider, Amount, PaymentDt, Bank, DeliveryCost, GoodsTotal, CustomFee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.Payment.Transaction, data.Payment.RequestID, data.Payment.Currency, data.Payment.Provider, data.Payment.Amount, data.Payment.PaymentDt, data.Payment.Bank, data.Payment.DeliveryCost, data.Payment.GoodsTotal, data.Payment.CustomFee)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *DataBase) insertItems(data user.UserData) {
	stmt, err := v.db.Prepare("INSERT INTO items (ChrtID, TrackNumber, Price, Rid, Name, Sale, Size, TotalPrice, NmID, Brand, Status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < len(data.Items); i++ {
		_, err = stmt.Exec(data.Items[0].ChrtID, data.Items[0].TrackNumber, data.Items[0].Price, data.Items[0].Rid, data.Items[0].Name, data.Items[0].Sale, data.Items[0].Size, data.Items[0].TotalPrice, data.Items[0].NmID, data.Items[0].Brand, data.Items[0].Status)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (v *DataBase) insertOrders(data user.UserData) {
	stmt, err := v.db.Prepare("INSERT INTO orders (OrderUID, TrackNumber, Entry, Locale, InternalSignature, CustomerID, DeliveryService, Shardkey, SmID, DateCreated, OofShard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < len(data.Items); i++ {
		_, err = stmt.Exec(data.OrderUID, data.TrackNumber, data.Entry, data.Locale, data.InternalSignature, data.CustomerID, data.DeliveryService, data.Shardkey, data.SmID, data.DateCreated, data.OofShard)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (v *DataBase) Listen() stan.Subscription {
	sc, _ := stan.Connect("test-cluster", "client_db", stan.NatsURL("nats://0.0.0.0:4222"))
	sub, _ := sc.Subscribe("database", func(msg *stan.Msg) {
		var task task.Task
		err := json.Unmarshal(msg.Data, &task)
		if err != nil {
			log.Fatal(err)
		}

		if task.UpdateDB {
			v.ReadFile(task.NameFile)
		}

		if task.Cash {
			v.Regenerate()
		}
	})
	return sub
}

func (v *DataBase) Regenerate() {
	// var data user.UserData
	var sub, _ = v.db.Exec("SELECT * FROM delivery")
	log.Print(sub)

	// v.db.Exec("SELECT * FROM payment")
	// v.db.Exec("SELECT * FROM items")
	// v.db.Exec("SELECT * FROM orders")
}
