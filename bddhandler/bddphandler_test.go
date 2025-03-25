package bddphandler

import (
    "testing"
    _ "github.com/mattn/go-sqlite3"
)

func TestInsertUser(t *testing.T) {
    db, err :=InitializeDB(":memory:")
    if err != nil {
        t.Fatalf("Could not initialize test database: %v", err) // Si jamais ça plante
    }
    defer CloseDB(db)

    err = createTables(db)
    if err != nil {
        t.Fatalf("Failed to create tables: %v", err)
    }

    _, err = InsertUser(db, "test_pelo", "test@gogole.zaza", "testmdp")
    if err != nil {
        t.Errorf("Failed to insert user: %v", err)
    }
}
