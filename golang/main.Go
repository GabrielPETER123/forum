package main

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/GabrielPETER123/Forum/database"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("index.html"))
  if r.Method == http.MethodGet {
  }else{
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }
    textPost := r.FormValue("textPost")
    fmt.Println(textPost)
  }
  
  tmpl.Execute(w, nil)
}

func main() {
  database.DataBase()
  
  http.HandleFunc("/", indexHandler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
  fmt.Println("(http://localhost:8080/)")
  http.ListenAndServe(":8080", nil)
}