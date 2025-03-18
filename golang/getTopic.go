package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func GetTopic(topicID int) Topic{
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	db.AutoMigrate(&Topic{})

	//* Lecture de la base de données
	var topic Topic
	if err := db.First(&topic, "ID = ?", topicID).Error; err != nil {
		fmt.Print("Le topic n'est pas présent dans la base de données\n")
	}

    topic.FormattedCreationDate = topic.CreatedAt.Format("02 January 2006 15:04")
    topic.FormattedUpdatedDate = topic.UpdatedAt.Format("02 January 2006 15:04")

	return topic
}

func GetAllTopics() []Topic {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	db.AutoMigrate(&Topic{})

	//* Copie des topics
	var topics []Topic
	db.Find(&topics)

	return topics
}