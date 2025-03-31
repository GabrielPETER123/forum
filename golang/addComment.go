package golang

import (
	// "fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func AddComment(comment Comment) {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Migation du schéma
	// fmt.Println("Migrating schema...")
	db.AutoMigrate(&Comment{})
	// fmt.Println("Schema migrated.")

	//* Création du Post
	// fmt.Println("Creating comment...")
	db.Create(&comment)
	// fmt.Println("Comment created.")
}