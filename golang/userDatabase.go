package golang

import (
	// "fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func AddUserInDataBase (userSend User){
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Migation du schéma
	// fmt.Println("Migrating schema...")
	db.AutoMigrate(&User{})
	// fmt.Println("Schema migrated.")

	//* Création du Post
	// fmt.Println("Creating user...")
	db.Create(&userSend)
	// fmt.Println("User created.")
}