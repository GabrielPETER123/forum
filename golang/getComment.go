package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func GetCommentByPostID(postID uint) []Comment {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Lecture de la base de données
	var comments []Comment

	//* Recherche des commentaires du post
	if err := db.Where("PostID = ?", postID).Find(&comments).Error; err != nil {
		fmt.Println("Comments not found:", err)
		return nil
	}

	//* Formatage des dates
	for i := range comments {
        comments[i].FormattedCreationDate = comments[i].CreatedAt.Format("02 January 2006 15:04")
        comments[i].FormattedUpdatedDate = comments[i].UpdatedAt.Format("02 January 2006 15:04")
    }

	return comments
}

func GetCommentsByUserID(userID int) []Comment {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")
	
	//* Lecture de la base de données
	var comments []Comment

	//* Recherche des commentaires de l'utilisateur
	if err := db.Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		fmt.Println("Comments not found:", err)
		return nil
	}

	//* Formatage des dates
	for i := range comments {
		comments[i].FormattedCreationDate = comments[i].CreatedAt.Format("02 January 2006 15:04")
		comments[i].FormattedUpdatedDate = comments[i].UpdatedAt.Format("02 January 2006 15:04")
	}

	return comments
}

func GetAllComments() []Comment {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Lecture de la base de données
	var comments []Comment

	//* Recherche de tous les commentaires
	if err := db.Find(&comments).Error; err != nil {
		fmt.Println("Comments not found:", err)
		return nil
	}

	//* Formatage des dates
	for i := range comments {
		comments[i].FormattedCreationDate = comments[i].CreatedAt.Format("02 January 2006 15:04")
		comments[i].FormattedUpdatedDate = comments[i].UpdatedAt.Format("02 January 2006 15:04")
	}

	return comments
}