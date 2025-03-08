package main

import (
  "fmt"
  "net/http"
  "html/template"
  "forum/golang"
)

//* Fonction qui gère la page d'accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/index.html"))
  if r.Method == http.MethodGet {
  }else{
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }

    title := r.FormValue("title")
    content := r.FormValue("content")

    fmt.Println(title, content)

    var postSend = golang.Post{}
    postSend.Title = title
    postSend.Text = content

    //* Écris dans la base de données si le titre ou le contenu n'est pas vide
    if title != "" || content != "" {
    fmt.Println("Starting database...")
    golang.PostDataBase(postSend)
    fmt.Println("Database ended.")
    }

  }
  tmpl.Execute(w, nil)
}

func main() {
  http.HandleFunc("/", indexHandler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("(http://localhost:8080/)")
	http.ListenAndServe(":8080", nil)
}