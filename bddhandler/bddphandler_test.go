package bddphandler_test

import (
    "testing"
    "bddphandler"
    _ "github.com/mattn/go-sqlite3"
)

func TestInsertUser(t *testing.T) {
    db, err := bddphandler.InitializeDB(":memory:")
    if err != nil {
        t.Fatalf("Could not initialize test database: %v", err)
    }
    defer bddphandler.CloseDB(db)

    err = bddphandler.createTables(db)
    if err != nil {
        t.Fatalf("Failed to create tables: %v", err)
    }

    _, err = bddphandler.InsertUser(db, "test_pelo", "test@gogole.zaza", "testmdp")
    if err != nil {
        t.Errorf("Failed to insert user: %v", err)
    }
}