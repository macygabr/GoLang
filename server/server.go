package server

import (
	"html/template"
	"io"
	"level0/database"
	"net/http"
	"os"
)

type Server struct {
	db *database.DataBase
}

func (s *Server) Start(db *database.DataBase) {
	s.db = db
	db.Connect()

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
			s.db.ReadFile()
		} else {
			tmpl, _ := template.ParseFiles("server/templates/pages/upload.html")
			tmpl.Execute(w, nil)
		}
	})
}
