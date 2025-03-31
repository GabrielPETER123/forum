package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func AddTopic(name string, description string, user User) string {
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Topic{})

	var topic Topic

	//! Enpêche la création d'un topic déjà existant
	if err := db.First(&topic, "Name = ?", name).Error; err == nil {
		fmt.Print("Le topic est déjà présent dans la base de données\n")
		return "Le topic est déjà présent dans la base de données"
	}

	fmt.Println("Creating topic...")
	db.Create(&Topic{Name: name, Description: description, User : user})
	fmt.Println("Topic created.")
	
	return "Topic created"
}