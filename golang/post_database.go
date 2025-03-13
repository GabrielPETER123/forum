package golang

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Post struct {
	gorm.Model
	Title  string
	UserID uint
	User User `gorm:"foreignKey:UserID"`
	Text   string
}



//* Fonction qui écrit dans la base de données
func AddPostInDataBase (postSend Post){
	fmt.Println("Opening database connection...")
	db, err := gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connection opened.")

	//* Migation du schéma
	fmt.Println("Migrating schema...")
	db.AutoMigrate(&Post{})
	fmt.Println("Schema migrated.")

	//* Lecture de la base de données
	//! TEST si le texte est déjà présent dans la base de données
	var post Post
	if err := db.First(&post, "Text = ?", postSend.Text).Error; err == nil {
		fmt.Print("Le texte est déjà présent dans la base de données\n")
		return
	}

	//* Création du Post
	fmt.Println("Creating post...")
	db.Create(&postSend)
	fmt.Println("Post created.")
}

	// // Read
	// fmt.Println("Reading post...")
	// var product Post
	// db.First(&product, 1) // find post with integer primary key
	// db.First(&product, "code = ?", "D42") // find post with code D42
	// fmt.Println("Post read.")

	// // Update - update post's price to 200
	// fmt.Println("Updating post price...")
	// db.Model(&product).Update("Title", "Bière blonde")
	// // Update - update multiple fields
	// db.Model(&product).Updates(Post{Text: "La bière c'est trop bien"}) // non-zero fields
	// fmt.Println("Post updated.")