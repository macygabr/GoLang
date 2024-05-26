package server

import (
	"fmt"
	"golang/database"
	"html/template"
	"io"
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

func Send() {
	fmt.Println("Send")
	sc, _ := stan.Connect("test-cluster", "client_server", stan.NatsURL("nats://0.0.0.0:4222"))
	defer sc.Close()
	message := []byte("Hello, NATS Streaming!!!")
	sc.Publish("parseFile", message)
}

func (s *Server) Start() {
	s.ListenHomePage()
	s.ListenUploadPage()
	http.ListenAndServe(":8080", nil)
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
