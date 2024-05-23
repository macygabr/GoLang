package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

type MyData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Connect() {
	conninfo := "user=postgres password=postgres host=127.0.0.1 sslmode=disable"
	db, err := sql.Open("postgres", conninfo)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS godb ( id integer, username varchar(255) )")

	if err != nil {
		log.Fatal(err)
	}
}

func ReadFile() {
	// Чтение содержимого JSON-файла
	jsonData, err := ioutil.ReadFile("./download/test.json")
	if err != nil {
		log.Fatal(err)
	}

	// Декодирование данных JSON в структуру или карту
	var data []MyData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Вывод длины данных JSON
	fmt.Println(len(jsonData))

	// Вывод декодированных данных
	fmt.Println(data)

	// Установка соединения с базой данных PostgreSQL
	db, err := sql.Open("postgres", "user=your_user password=your_password dbname=your_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Подготовка SQL-запроса для вставки данных в таблицу
	stmt, err := db.Prepare("INSERT INTO godb (id, name) VALUES ($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Выполнение SQL-запроса для каждого элемента данных
	for _, item := range data {
		_, err = stmt.Exec(item.ID, item.Name)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Данные успешно вставлены.")
}
