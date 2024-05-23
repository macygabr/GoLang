package main

import (
	"fmt"
	"html/template"
	"io"
	"level0/database"
	"level0/model"
	"net/http"
	"os"
)

func main() {
	database.Connect()
	ListenHomePage()
	ListenUploadPage()
	http.ListenAndServe(":8080", nil)
}

func ListenHomePage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := model.ViewData{
			Title:   "home",
			Message: "nurlan",
		}
		tmpl, _ := template.ParseFiles("templates/pages/home.html")
		tmpl.Execute(w, data)
	})
}

func ListenUploadPage() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(r)
		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			savePath := "download/" + header.Filename
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
			// fmt.Println("Ok")
			tmpl, _ := template.ParseFiles("templates/pages/upload.html")
			tmpl.Execute(w, nil)
			fmt.Println("Ok")
			database.ReadFile()
		} else {
			tmpl, _ := template.ParseFiles("templates/pages/upload.html")
			tmpl.Execute(w, nil)
		}
	})
}
