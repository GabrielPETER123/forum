package golang

import (
	// "fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

//* Fonction qui met à jour un post dans la base de données
func UpdatePost(postSend Post){
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		// panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Migation du schéma
	// fmt.Println("Migrating schema...")
	db.AutoMigrate(&Post{})
	// fmt.Println("Schema migrated.")

	//* Lecture de la base de données
	var post Post

	//* Recherche du post à mettre à jour
	if err := db.First(&post, postSend.ID).Error; err != nil {
		// fmt.Println("Post not found:", err)
		return
	}
	//* Mise à jour du Post
	// fmt.Println("Updating post...")
	db.Model(&post).Updates(postSend)
	// fmt.Println("Post updated.")
}