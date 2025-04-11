package golang

import (
	"golang.org/x/crypto/bcrypt"
)

// * création de l'utilisateur admin
func CreateAdminUser() {
	if CheckUser("admin") {
		return
	}
	password := "admin"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user := User{
		Id:       "admin",
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "",
		Admin:    true,
	}
	AddUserInDataBase(user)
}
