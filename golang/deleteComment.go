package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func DeleteComment(commentID uint) {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Suppression du commentaire
	if err := db.Delete(&Comment{}, commentID).Error; err != nil {
		fmt.Println("Comment not found:", err)
		return
	}
}