package main

import (
  "fmt"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)


type Post struct {
	gorm.Model
	Title  string
	User string
	UserID uint
	Text   string
}

func DataBase () {
	fmt.Println("Connexion à la base de données...")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
	  panic("Erreur lors de la connexion à la base de données.")
	}
	fmt.Println("Connexion réussie.")
  
	// Migrate the schema
	fmt.Println("Migration du schéma...")
	db.AutoMigrate(&Post{})
	fmt.Println("Migration réussie.")
  
	// Create
	fmt.Println("Création d'un post...")
	db.Create(&Post{Title: "Bière", User: "Rodrigo", UserID: 1, Text: "Bière blonde"})
	fmt.Println("Post créé.")
  
	// Read
	fmt.Println("Lecture du post...")
	var product Post
	db.First(&product, 1) // find post with integer primary key
	db.First(&product, "code = ?", "D42") // find post with code D42
	fmt.Println("Post lu.")
  
	// Update - update post's price to 200
	fmt.Println("Mise à jour du post...")
	db.Model(&product).Update("Title", "Bière blonde")
	// Update - update multiple fields
	db.Model(&product).Updates(Post{Text: "La bière c'est trop bien"}) // non-zero fields
	fmt.Println("Post mis à jour.")
}