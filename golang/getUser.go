package golang

import (
	// "fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func GetUserByID(userID string) User {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Copie de l'utilisateur
	var user User
	db.Where("id = ?", userID).Find(&user)

	user.FormattedCreationDate = user.CreatedAt.Format("02 January 2006 15:04")
    user.FormattedUpdatedDate = user.UpdatedAt.Format("02 January 2006 15:04")

	return user
}

func GetAllUsers() []User {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Copie des utilisateurs
	var users []User
	db.Find(&users)
	return users
}