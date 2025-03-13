package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func CheckUser(nameOrMail string, password string) bool{
	fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connection opened.")

	//* Migation du schéma
	fmt.Println("Migrating schema...")
	db.AutoMigrate(&User{})
	fmt.Println("Schema migrated.")

	//* Lecture de la base de données
	var user User
	if err := db.First(&user, "Username = ? AND Password = ?", nameOrMail, password).Error; err == nil {
		fmt.Print("L'utilisateur est présent dans la base de données\n")
		return true
	}
	fmt.Print("L'utilisateur n'est pas présent dans la base de données\n")
	return false
}