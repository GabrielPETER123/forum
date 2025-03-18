package golang

import (
    "gorm.io/gorm"
)

//* Définition des structures de données

type User struct {
    gorm.Model
    Username string
    Password string
    Email    string
    Admin    bool
    Posts    []Post
}

type Post struct {
    gorm.Model
    Title                 string
    Text                  string
    UserID                uint
    User                  User
    FormattedCreationDate string
    FormattedUpdatedDate  string
    TotalUp               int
    TotalDown             int
    Votes                 []Vote `gorm:"foreignKey:PostID"`
}

type Vote struct {
    ID     uint `gorm:"primaryKey"`
    PostID uint
    UserID uint
    Up    int
    Down  int
}

