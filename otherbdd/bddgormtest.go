package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:100"`
	Email    string `gorm:"size:100;unique_index"`
	Articles []Article
}

type Article struct {
	ID     uint   `gorm:"primary_key"`
	Title  string `gorm:"size:200"`
	Body   string `gorm:"size:1000"`
	UserID uint
	User   User
}

func main() {
	db, err := gorm.Open("sqlite3", "gorm_example.db")
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Article{})

	user := User{Name: "John Doe", Email: "johndoe@example.com"}
	db.Create(&user)

	article1 := Article{Title: "Premier article", Body: "Ceci est le premier article de John."}
	article2 := Article{Title: "Deuxième article", Body: "Ceci est le deuxième article de John."}
	db.Model(&user).Association("Articles").Append(&article1, &article2)

	var retrievedUser User
	if err := db.Preload("Articles").First(&retrievedUser, "email = ?", "eminem@example.com").Error; err != nil {
		log.Fatalf("Erreur de récupération de l'utilisateur: %v", err)
	}
	fmt.Printf("Utilisateur récupéré: %v\n", retrievedUser.Name)
	for _, article := range retrievedUser.Articles {
		fmt.Printf("Article: %v - %v\n", article.Title, article.Body)
	}

	db.Model(&retrievedUser).Update("Name", "John Updated")

	var articleToDelete Article
	db.First(&articleToDelete, "title = ?", "Premier article")
	db.Delete(&articleToDelete)

	var remainingArticles []Article
	db.Where("user_id = ?", retrievedUser.ID).Find(&remainingArticles)
	fmt.Printf("Articles restants après suppression: %v\n", len(remainingArticles))
}
