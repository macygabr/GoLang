package database

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"golang/model/item"
	"golang/model/task"
	"golang/model/user"
	"log"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
)

type DataBase struct {
	db   *sql.DB
	user user.UserData
}

func (v *DataBase) Connect() stan.Subscription {
	var sub = v.Listen()
	conninfo := "user=postgres password=postgres host=127.0.0.1 sslmode=disable"
	db, err := sql.Open("postgres", conninfo)

	if err != nil {
		log.Fatal(err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS delivery (ID serial primary key, Name varchar(255), Phone varchar(255), Zip varchar(255), City varchar(255), Address varchar(255), Region varchar(255), Email varchar(255) )")
	db.Exec("CREATE TABLE IF NOT EXISTS payment (ID serial primary key, Transaction varchar(255) , RequestID varchar(255), Currency varchar(255), Provider varchar(255), Amount INTEGER, PaymentDt INTEGER, Bank varchar(255), DeliveryCost INTEGER, GoodsTotal INTEGER, CustomFee INTEGER)")
	db.Exec("CREATE TABLE IF NOT EXISTS items (ID serial primary key, order_id varchar(255), ChrtID INTEGER, TrackNumber varchar(255), Price INTEGER, Rid varchar(255), Name varchar(255), Sale INTEGER, Size varchar(255), TotalPrice INTEGER, NmID INTEGER, Brand varchar(255), Status INTEGER)")
	db.Exec("CREATE TABLE IF NOT EXISTS orders(delivery_id INTEGER, payment_id INTEGER, items_id INTEGER, OrderUID varchar(255), TrackNumber varchar(255), Entry varchar(255), Locale varchar(255), InternalSignature varchar(255), CustomerID varchar(255), DeliveryService varchar(255), Shardkey varchar(255), SmID INTEGER, DateCreated varchar(255), OofShard varchar(255))")

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

	err = json.Unmarshal(jsonData, &v.user)
	if err != nil {
		log.Fatal(err)
	}
	v.insertDelivery(v.user)
	v.insertPayment(v.user)
	v.insertItems(v.user)
	v.insertOrders(v.user)
	v.Regenerate()
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
	stmt, err := v.db.Prepare("INSERT INTO items (order_id, ChrtID, TrackNumber, Price, Rid, Name, Sale, Size, TotalPrice, NmID, Brand, Status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < len(data.Items); i++ {
		_, err = stmt.Exec(data.OrderUID, data.Items[0].ChrtID, data.Items[0].TrackNumber, data.Items[0].Price, data.Items[0].Rid, data.Items[0].Name, data.Items[0].Sale, data.Items[0].Size, data.Items[0].TotalPrice, data.Items[0].NmID, data.Items[0].Brand, data.Items[0].Status)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (v *DataBase) insertOrders(data user.UserData) {
	stmt, err := v.db.Prepare("INSERT INTO orders (delivery_id, payment_id, items_id, OrderUID, TrackNumber, Entry, Locale, InternalSignature, CustomerID, DeliveryService, Shardkey, SmID, DateCreated, OofShard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var payment_id int
	v.db.QueryRow("SELECT MAX(id) FROM payment").Scan(&payment_id)

	var delivery_id int
	v.db.QueryRow("SELECT MAX(id) FROM delivery").Scan(&delivery_id)

	var items_id int
	v.db.QueryRow("SELECT MAX(id) FROM items").Scan(&items_id)

	for i := 1; i <= len(data.Items); i++ {
		_, err = stmt.Exec(delivery_id, payment_id, items_id-len(data.Items)+i, data.OrderUID, data.TrackNumber, data.Entry, data.Locale, data.InternalSignature, data.CustomerID, data.DeliveryService, data.Shardkey, data.SmID, data.DateCreated, data.OofShard)
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
	sc, _ := stan.Connect("test-cluster", "db_send", stan.NatsURL("nats://0.0.0.0:4222"))
	defer sc.Close()

	task := new(task.Task)
	task.SetUpdateDB(true)

	rows, _ := v.db.Query("select orderuid, orders.tracknumber, orders.entry, delivery.name, delivery.phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email,  payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.paymentdt, payment.Bank, payment.deliverycost, payment.GoodsTotal, payment.CustomFee, orders.locale, orders.CustomerID, orders.DeliveryService, orders.Shardkey, orders.SmID, orders.DateCreated, orders.oofshard FROM orders JOIN delivery ON orders.delivery_id = delivery.id JOIN payment ON orders.payment_id = payment.id")
	defer rows.Close()

	for rows.Next() {
		user := new(user.UserData)
		err := rows.Scan(&user.OrderUID, &user.TrackNumber, &user.Entry,
			&user.Delivery.Name, &user.Delivery.Phone, &user.Delivery.Zip, &user.Delivery.City, &user.Delivery.Address, &user.Delivery.Region, &user.Delivery.Email,
			&user.Payment.RequestID, &user.Payment.Currency, &user.Payment.Provider, &user.Payment.Amount, &user.Payment.PaymentDt, &user.Payment.Bank, &user.Payment.DeliveryCost, &user.Payment.GoodsTotal, &user.Payment.CustomFee,
			&user.Locale, &user.CustomerID, &user.DeliveryService, &user.Shardkey, &user.SmID, &user.DateCreated, &user.OofShard)

		items, _ := v.db.Query("SELECT chrtid, tracknumber, price, rid, name, sale, size, totalprice, nmid, brand, status FROM items WHERE order_id = $1", user.OrderUID)
		defer items.Close()
		for items.Next() {

			item := new(item.Items)
			// log.Print(items)
			err := items.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				log.Println(err)
				continue
			}
			user.Items = append(user.Items, *item)
		}

		if err != nil {
			log.Println(err)
			continue
		}
		task.SetUserData(*user)
		// log.Println(task)

		message, err := json.Marshal(task)
		if err != nil {
			log.Fatal(err)
		}
		sc.Publish("cash", message)
	}
}
