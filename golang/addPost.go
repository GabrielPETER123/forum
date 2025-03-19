package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

//* Fonction qui écrit dans la base de données
func AddPostInDataBase (postSend Post){
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Migation du schéma
	// fmt.Println("Migrating schema...")
	db.AutoMigrate(&Post{})
	// fmt.Println("Schema migrated.")

	//* Lecture de la base de données
	//! TEST si le texte est déjà présent dans la base de données
	var post Post
	if err := db.First(&post, "Text = ?", postSend.Text).Error; err == nil {
		fmt.Print("Le texte est déjà présent dans la base de données\n")
		return
	}

	//* Création du Post
	// fmt.Println("Creating post...")
	db.Create(&postSend)
	// fmt.Println("Post created.")
}
