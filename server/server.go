package server

import (
	"encoding/json"
	"golang/database"
	"golang/model/task"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/stan.go"
)

type Server struct {
	db *database.DataBase
	sc stan.Conn
}

func NewServer(db *database.DataBase, arg stan.Conn) *Server {
	return &Server{db, arg}
}

func (s *Server) Start() stan.Subscription {
	var sub = s.Listen()
	s.ListenHomePage()
	s.ListenUploadPage()
	s.ListenOrderPage()

	http.ListenAndServe(":8080", nil)
	return sub
}

func (s *Server) ListenHomePage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("server/templates/pages/home.html")
		tmpl.Execute(w, nil)
	})
}

func (s *Server) ListenUploadPage() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			savePath := "server/download/" + header.Filename
			dst, err := os.Create(savePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl, _ := template.ParseFiles("server/templates/pages/upload.html")
			tmpl.Execute(w, nil)
			Send()
		} else {
			tmpl, _ := template.ParseFiles("server/templates/pages/upload.html")
			tmpl.Execute(w, nil)
		}
	})
}

func (s *Server) ListenOrderPage() {
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		sc, _ := stan.Connect("test-cluster", "client_server_send", stan.NatsURL("nats://0.0.0.0:4222"))
		defer sc.Close()

		task := new(task.Task)
		task.SetCash(true)

		message, err := json.Marshal(task)
		if err != nil {
			log.Fatal(err)
		}
		sc.Publish("cash", message)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	})
}

func Send() {
	sc, _ := stan.Connect("test-cluster", "client_server_send", stan.NatsURL("nats://0.0.0.0:4222"))
	defer sc.Close()

	task := new(task.Task)
	task.SetUpdateDB(true)

	message, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	sc.Publish("Task", message)
}

func (s *Server) Listen() stan.Subscription {
	sc, _ := stan.Connect("test-cluster", "client_server_listen", stan.NatsURL("nats://0.0.0.0:4222"))
	sub, _ := sc.Subscribe("server", func(msg *stan.Msg) {
		var task task.Task
		err := json.Unmarshal(msg.Data, &task)
		if err != nil {
			log.Fatal(err)
		}

		if task.UpdateDB {
			log.Println(task)
		}

		if task.Cash {
			log.Println(task)
		}
	})
	return sub
}
