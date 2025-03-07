package bddphandler

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

// Struct pour user (gabriel nul)
type User struct {
    ID        int
    Username  string
    Email     string
    CreatedAt time.Time
}

// Struct pour post forum (rodrigo nul)
type Post struct {
    ID        int
    UserID    int
    Title     string
    Content   string
    CreatedAt time.Time
}

func InitializeDB(dbName string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }

    createUsersTable := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

    // Table posts
    createPostsTable := `CREATE TABLE IF NOT EXISTS posts (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

    createCommentsTable := `CREATE TABLE IF NOT EXISTS comments (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

    // exéc de création des tables (faut pas être con c'est marqué exéc c'est juste moi)
    for _, tableSQL := range []string{createUsersTable, createPostsTable, createCommentsTable} {
        if _, err := db.Exec(tableSQL); err != nil {
            return nil, fmt.Errorf("failed to create table: %v", err)
        }
    }

    return db, nil
}

func InsertUser(db *sql.DB, username, email, passwordHash string) (int64, error) {
    insertUserSQL := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
    result, err := db.Exec(insertUserSQL, username, email, passwordHash)
    if err != nil {
        return 0, fmt.Errorf("failed to insert user: %v", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert id: %v", err)
    }
    
    return id, nil
}

func CreatePost(db *sql.DB, userID int, title, content string) (int64, error) {
    insertPostSQL := `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
    result, err := db.Exec(insertPostSQL, userID, title, content)
    if err != nil {
        return 0, fmt.Errorf("failed to create post: %v", err)
    }
    
    return result.LastInsertId()
}

func GetUserPosts(db *sql.DB, userID int) ([]Post, error) {
    query := `SELECT id, title, content, created_at FROM posts WHERE user_id = ?`
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user posts: %v", err)
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan post: %v", err)
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func AddComment(db *sql.DB, postID, userID int, content string) error {
    insertCommentSQL := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
    _, err := db.Exec(insertCommentSQL, postID, userID, content)
    if err != nil {
        return fmt.Errorf("failed to add comment: %v", err)
    }
    return nil
}

func DeletePost(db *sql.DB, postID, userID int) error {
    // ça vérif si l'user c'est bien le créateur du post
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }

    // là ça delete les coms
    _, err = tx.Exec(`DELETE FROM comments WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to delete comments: %v", err)
    }

    // là ça delete le post
    result, err := tx.Exec(`DELETE FROM posts WHERE id = ? AND user_id = ?`, postID, userID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to delete post: %v", err)
    }

    affected, err := result.RowsAffected()
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to get affected rows: %v", err)
    }

    if affected == 0 {
        tx.Rollback()
        return fmt.Errorf("post not found or user not authorized")
    }

    return tx.Commit()
}

func CloseDB(db *sql.DB) {
    if err := db.Close(); err != nil {
        log.Printf("Error closing database: %v", err)
    }
}