package golang

import (
    "gorm.io/gorm"
)

//* Définition des structures de données

type User struct {
    gorm.Model
    Username string `gorm:"size:255;not null;unique"`
    Email    string `gorm:"size:255;not null;unique"`
    Password string `gorm:"size:255;not null"`
    Admin    bool
    Posts    []Post
    TotalVote uint
    TotalPost uint
}

type Post struct {
    gorm.Model
    Title                 string
    Text                  string `gorm:"size:255;not null"`
    UserID                uint
    User                  User
    FormattedCreationDate string `gorm:"-"`
    FormattedUpdatedDate  string `gorm:"-"`
    TotalUp               uint
    TotalDown             uint
    Votes                 []Vote `gorm:"foreignKey:PostID"`
    Comments              []Post `gorm:"foreignKey:ParentID"`
    ParentID              uint   `gorm:"default:null"`
    TopicID               uint   `gorm:"not null"`
    IsLoggedIn            bool   `gorm:"-"`
    UserConnectedID       uint   `gorm:"-"`
}

type Vote struct {
    ID     uint `gorm:"primaryKey"`
    PostID uint
    UserID uint
    Up    int
    Down  int
}

type Topic struct {
    gorm.Model
    Name        string `gorm:"size:255;not null;unique"`
    Description string `gorm:"size:255"`
    UserID      uint   `gorm:"not null"`
    User        User   `gorm:"foreignKey:UserID"`
    FormattedCreationDate string `gorm:"-"`
    FormattedUpdatedDate  string `gorm:"-"`
    Posts       []Post `gorm:"foreignKey:TopicID"`
}
