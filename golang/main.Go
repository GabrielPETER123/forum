package main

import (
  "fmt"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "net/http"
  "html/template"
)

type Post struct {
  gorm.Model
  Title  string
  User string
  UserID uint
  Text   string
}


func dataBase (){
  fmt.Println("Opening database connection...")
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }
  fmt.Println("Database connection opened.")

  // Migrate the schema
  fmt.Println("Migrating schema...")
  db.AutoMigrate(&Post{})
  fmt.Println("Schema migrated.")

  // Create
  fmt.Println("Creating post...")
  db.Create(&Post{Title: "Bière", User: "Rodrigo", UserID: 1, Text: "Bière blonde"})
  fmt.Println("Post created.")

  // Read
  fmt.Println("Reading post...")
  var product Post
  db.First(&product, 1) // find post with integer primary key
  db.First(&product, "code = ?", "D42") // find post with code D42
  fmt.Println("Post read.")

  // Update - update post's price to 200
  fmt.Println("Updating post price...")
  db.Model(&product).Update("Title", "Bière blonde")
  // Update - update multiple fields
  db.Model(&product).Updates(Post{Text: "La bière c'est trop bien"}) // non-zero fields
  fmt.Println("Post updated.")
}

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
  http.HandleFunc("/", indexHandler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("(http://localhost:8080/)")
	http.ListenAndServe(":8080", nil)
}