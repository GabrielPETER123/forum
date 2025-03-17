package golang

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetPostsByUserID(userID int) []Post {
    fmt.Println("Opening database connection...")
    db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    fmt.Println("Database connection opened.")

    //* Copie des posts de l'utilisateur
    var posts []Post
    db.Where("user_id = ?", userID).Find(&posts)
    return posts
}