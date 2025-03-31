package golang

import (
    // "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "golang.org/x/crypto/bcrypt"
)

func CheckUserPassword(nameOrMail, password string) (User, bool) {
    // fmt.Println("Opening database connection...")
    db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    // fmt.Println("Database connection opened.")

    //* Regarde si l'utilisateur est présent dans la base de données
    var user User
    if err := db.Where("username = ? OR email = ?", nameOrMail, nameOrMail).First(&user).Error; err == nil {
        //* Compare le mot de passe avec le hash password
        err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
        if err == nil {
            return user, true
        }
    }
    return User{}, false
}