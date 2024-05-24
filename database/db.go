package database

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"level0/model"
	"log"

	_ "github.com/lib/pq"
)

type DataBase struct {
	db *sql.DB
}

func (v *DataBase) Connect() {
	conninfo := "user=postgres password=postgres host=127.0.0.1 sslmode=disable"
	db, err := sql.Open("postgres", conninfo)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS godb ( OrderUID varchar(255) , Name varchar(255) )")

	if err != nil {
		log.Fatal(err)
	}

	v.db = db
}

func (v *DataBase) ReadFile() {
	jsonData, err := ioutil.ReadFile("server/download/model.json")
	if err != nil {
		log.Fatal(err)
	}

	if !json.Valid(jsonData) {
		return
	}

	var data model.UserData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := v.db.Prepare("INSERT INTO godb (OrderUID, Name) VALUES ($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.OrderUID, data.Delivery.Name)
	if err != nil {
		log.Fatal(err)
	}
}
