package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func DeletePost(id int) {
	fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connection opened.")

	//* Suppression du post
	fmt.Println("Deleting post...")
	db.Where("ID = ?", id).Delete(&Post{})
	fmt.Println("Post deleted.")
}