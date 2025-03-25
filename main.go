package main

import (
	"fmt"
	"log"
	"bddphandler"
)

func main() {
	db, err := bddphandler.InitDB("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Base de données initialisée avec succès.")
}
