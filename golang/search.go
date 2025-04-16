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

func SearchUsersByUsername(username string) []User {
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var users []User

	// Rechercher dans les utilisateurs
	db.Where("username LIKE ?", "%"+username+"%").Find(&users)

	for i := range users {
		users[i].FormattedCreationDate = users[i].CreatedAt.Format("02 January 2006 15:04")
		users[i].FormattedUpdatedDate = users[i].UpdatedAt.Format("02 January 2006 15:04")
	}

	return users
}

func SearchPostsByTitle(title string) []Post {
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var posts []Post

	// Rechercher dans les posts
	db.Where("title LIKE ?", "%"+title+"%").Find(&posts)

	for i := range posts {
		posts[i].FormattedCreationDate = posts[i].CreatedAt.Format("02 January 2006 15:04")
		posts[i].FormattedUpdatedDate = posts[i].UpdatedAt.Format("02 January 2006 15:04")
		posts[i].TotalUp, posts[i].TotalDown = Totals(posts[i].ID)
	}

	return posts
}

func SearchTopicsByName(name string) []Topic {
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var topics []Topic

	// Rechercher dans les topics
	db.Where("name LIKE ?", "%"+name+"%").Find(&topics)

	for i := range topics {
		topics[i].FormattedCreationDate = topics[i].CreatedAt.Format("02 January 2006 15:04")
		topics[i].FormattedUpdatedDate = topics[i].UpdatedAt.Format("02 January 2006 15:04")
	}

	return topics
}

func SearchCommentsByText(text string) []Comment {
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var comments []Comment

	// Rechercher dans les commentaires
	db.Where("text LIKE ?", "%"+text+"%").Find(&comments)

	for i := range comments {
		comments[i].FormattedCreationDate = comments[i].CreatedAt.Format("02 January 2006 15:04")
		comments[i].FormattedUpdatedDate = comments[i].UpdatedAt.Format("02 January 2006 15:04")
	}

	return comments
}