package golang

import (
	"gorm.io/gorm"
)

//* Définition des structures de données

type User struct {
	gorm.Model
	Id                    string `gorm:"primaryKey"`
	Username              string `gorm:"size:255;not null;unique"`
	Email                 string `gorm:"size:255;not null;unique"`
	Password              string `gorm:"size:255;not null"`
	Admin                 bool
	Posts                 []Post
	FormattedCreationDate string `gorm:"-"`
	FormattedUpdatedDate  string `gorm:"-"`
	TotalVote             uint   `gorm:"-"`
	TotalPost             uint   `gorm:"-"`
	TotalComment          uint   `gorm:"-"`
}

type Post struct {
	gorm.Model
	Title                 string `gorm:"size:255;not null"`
	Text                  string `gorm:"not null"`
	UserID                string `gorm:"column:user_id"`
	User                  User   `gorm:"foreignKey:UserID;references:Id"`
	FormattedCreationDate string `gorm:"-"`
	FormattedUpdatedDate  string `gorm:"-"`
	TotalUp               uint
	TotalDown             uint
	Votes                 []Vote    `gorm:"foreignKey:PostID"`
	Comments              []Comment `gorm:"foreignKey:PostID"`
	TopicID               uint      `gorm:"not null"`
	IsLoggedIn            bool      `gorm:"-"`
	UserConnectedID       string    `gorm:"-"`
}

type Vote struct {
	ID     uint `gorm:"primaryKey"`
	PostID uint
	UserID string `gorm:"column:user_id"`
	User   User   `gorm:"foreignKey:UserID;references:Id"`
	Up     int
	Down   int
}

type Topic struct {
	gorm.Model
	Name                  string `gorm:"size:255;not null;unique"`
	Description           string `gorm:"size:255"`
	UserID                string `gorm:"column:user_id"`
	User                  User   `gorm:"foreignKey:UserID;references:Id"`
	FormattedCreationDate string `gorm:"-"`
	FormattedUpdatedDate  string `gorm:"-"`
	Posts                 []Post `gorm:"foreignKey:TopicID"`
}

type Comment struct {
	gorm.Model
	Text                  string `gorm:"not null"`
	UserID                string `gorm:"column:user_id"`
	User                  User   `gorm:"foreignKey:UserID;references:Id"`
	PostID                uint
	Post                  Post
	IsLoggedIn            bool   `gorm:"-"`
	UserConnectedID       string `gorm:"-"`
	FormattedCreationDate string `gorm:"-"`
	FormattedUpdatedDate  string `gorm:"-"`
}
