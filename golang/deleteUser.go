package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func DeleteUser(nameOrMail string) {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Suppression de l'utilisateur
	fmt.Println("Deleting user...")
	db.Where("Username = ?", nameOrMail).Delete(&User{})
	fmt.Println("User deleted.")
}
