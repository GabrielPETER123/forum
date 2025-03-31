package golang

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SearchUserPostTopic(search string) ([]User, []Post, []Topic) {
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
	}

	var users []User
	var posts []Post
	var topics []Topic

	// Rechercher dans les utilisateurs
	db.Where("username LIKE ?", "%"+search+"%").Find(&users)

	// Rechercher dans les posts
	db.Where("title LIKE ?", "%"+search+"%").Find(&posts)

	// Rechercher dans les topics
	db.Where("name LIKE ?", "%"+search+"%").Find(&topics)

	for i := range posts {
		posts[i].FormattedCreationDate = posts[i].CreatedAt.Format("02 January 2006 15:04")
		posts[i].FormattedUpdatedDate = posts[i].UpdatedAt.Format("02 January 2006 15:04")
		posts[i].TotalUp, posts[i].TotalDown = Totals(posts[i].ID)
	}

	return users, posts, topics
}