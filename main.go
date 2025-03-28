package main

import (
	"fmt"
	"log"
	"bddhandler"
)

func main() {
	db, err := bddphandler.InitDB("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Base de données initialisée avec succès.")
}
