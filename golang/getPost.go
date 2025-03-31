package golang

import (
	// "fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetPostsByUserID(userID int) []Post {
    // fmt.Println("Opening database connection...")
    db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    // fmt.Println("Database connection opened.")

    //* Copie des posts de l'utilisateur
    var posts []Post
    db.Where("user_id = ?", userID).Find(&posts)

    for i := range posts {
        posts[i].FormattedCreationDate = posts[i].CreatedAt.Format("02 January 2006 15:04")
        posts[i].FormattedUpdatedDate = posts[i].UpdatedAt.Format("02 January 2006 15:04")
    }
	for i := range posts {
		posts[i].TotalUp, posts[i].TotalDown = Totals(posts[i].ID)
	}
	
    return posts
}

func GetPostByPostID(postID int) Post {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Copie du post
	var post Post
	db.Preload("User").First(&post, postID)

    post.FormattedCreationDate = post.CreatedAt.Format("02 January 2006 15:04")
    post.FormattedUpdatedDate = post.UpdatedAt.Format("02 January 2006 15:04")

	post.TotalUp, post.TotalDown = Totals(post.ID)
	
	return post
}

func GetAllPosts() []Post {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Copie des posts
	var posts []Post
	db.Preload("User").Find(&posts)

    for i := range posts {
        posts[i].FormattedCreationDate = posts[i].CreatedAt.Format("02 January 2006 15:04")
        posts[i].FormattedUpdatedDate = posts[i].UpdatedAt.Format("02 January 2006 15:04")
    }
	for i := range posts {
		posts[i].TotalUp, posts[i].TotalDown = Totals(posts[i].ID)
	}
	
	return posts
}

func GetPostsByTopicID(topicID int) []Post {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	//* Copie des posts du topic
	var posts []Post
	db.Preload("Comments.User").Where("topic_id = ?", topicID).Find(&posts)

	for i := range posts {
        posts[i].User = GetUserByID(int(posts[i].UserID))
    }

	for i := range posts {
        posts[i].FormattedCreationDate = posts[i].CreatedAt.Format("02 January 2006 15:04")
        posts[i].FormattedUpdatedDate = posts[i].UpdatedAt.Format("02 January 2006 15:04")
		for j := range posts[i].Comments {
			posts[i].Comments[j].FormattedCreationDate = posts[i].Comments[j].CreatedAt.Format("02 January 2006 15:04")
			posts[i].Comments[j].FormattedUpdatedDate = posts[i].Comments[j].UpdatedAt.Format("02 January 2006 15:04")
		}
	}
	
	return posts
}

