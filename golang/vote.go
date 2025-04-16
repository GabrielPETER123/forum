package golang

import (
	// "fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

//* Gestion du vote
func Votes(postID uint, userID string, vote string) {
	// fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// fmt.Println("Database connection opened.")

	var post Post
	db.First(&post, postID)

	//* Migation du schéma
	db.AutoMigrate(&Vote{})
	
	//* Vérifie si le vote existe
	var existingVote Vote
	result := db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingVote)

	if result.Error != nil {
		// fmt.Println("Creating vote...")
		if vote == "up" {
			db.Create(&Vote{PostID: postID, UserID: userID, Up: 1, Down: 0})
		} else {
			db.Create(&Vote{PostID: postID, UserID: userID, Up: 0, Down: 1})
		}
		// fmt.Println("Vote created.")
	} else {
		// fmt.Println("Updating vote...")
		if vote == "up" {
			if existingVote.Up == 0 {
				db.Model(&existingVote).Update("up", 1)
				db.Model(&existingVote).Update("down", 0)
			} else {
				db.Model(&existingVote).Update("up", 0)
			}
		} 
		if vote == "down" {
			if existingVote.Down == 0 {
				db.Model(&existingVote).Update("down", 1)
				db.Model(&existingVote).Update("up", 0)
			} else {
				db.Model(&existingVote).Update("down", 0)
			}
		}
		// fmt.Println("Vote updated.")
	}
}

func Totals(postID uint) (uint, uint) {
    db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    var upVotes int64
    var downVotes int64

    db.Model(&Vote{}).Where("post_id = ? AND up = 1", postID).Count(&upVotes)
    db.Model(&Vote{}).Where("post_id = ? AND down = 1", postID).Count(&downVotes)

    return uint(upVotes), uint(downVotes)
}

func TotalVotes(userID string) uint {
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var totalVotes int64

	//*Compte le nombre de votes de l'utilisateur
	db.Model(&Vote{}).Where("user_id = ?", userID).Count(&totalVotes)

	return uint(totalVotes)
}