package bddphandler

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/bcrypt"
)

// **Infos de base d'un user
type User struct {
    ID        int       // **ID unique de l'user
    Username  string    // **Pseudo
    Email     string    // **Adresse email
    CreatedAt time.Time // **Date de création de l'user
}

// **Infos de base pour un post
type Post struct {
    ID        int       // **ID unique du post
    UserID    int       // **ID de l'user ayant créé le post
    Title     string    // **Titre du post
    Content   string    // **Contenu du post
    CreatedAt time.Time // **Date de création du post
}

// **Initialise la BDD en créant le fichier (si besoin) et en créant les tables.
func InitializeDB(dbName string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbName) // **Connexion à la BDD SQLite
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }

    // **Crée les tables nécessaires si elles n'existent pas
    if err := createTables(db); err != nil {
        return nil, err
    }

    return db, nil // **Retourne la connexion à la base de données.
}

// **createTables crée les tables pour les utilisateurs, les posts, et les commentaires.
func createTables(db *sql.DB) error {
    createUsersTable := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

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

    // **Exécute chaque instruction de création de table.
    for _, tableSQL := range []string{createUsersTable, createPostsTable, createCommentsTable} {
        if _, err := db.Exec(tableSQL); err != nil {
            return fmt.Errorf("failed to create table: %v", err)
        }
    }
    return nil
}

// **InsertUser insère un nouvel utilisateur dans la base de données.
func InsertUser(db *sql.DB, username, email, password string) (int64, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // **Hash du mot de passe.
    if err != nil {
        return 0, fmt.Errorf("failed to hash password: %v", err)
    }

    insertUserSQL := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
    result, err := db.Exec(insertUserSQL, username, email, string(hashedPassword)) // **Insertion de l'utilisateur dans la BDD.
    if err != nil {
        return 0, fmt.Errorf("failed to insert user: %v", err)
    }

    id, err := result.LastInsertId() // **Récupère l'ID de l'user inséré.
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert id: %v", err)
    }

    return id, nil // **Retourne l'ID de l'utilisateur.
}

// **CreatePost crée un nouveau post dans la base de données.
func CreatePost(db *sql.DB, userID int, title, content string) (int64, error) {
    insertPostSQL := `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
    result, err := db.Exec(insertPostSQL, userID, title, content) // **Insertion du post dans la BDD.
    if err != nil {
        return 0, fmt.Errorf("failed to create post: %v", err)
    }

    return result.LastInsertId() // **Retourne l'ID du post créé.
}

// **GetUserPosts récupère les posts d'un user spécifique.
func GetUserPosts(db *sql.DB, userID int) ([]Post, error) {
    query := `SELECT id, title, content, created_at FROM posts WHERE user_id = ?`
    rows, err := db.Query(query, userID) // **Exécution de la requête pour récupérer les posts.
    if err != nil {
        return nil, fmt.Errorf("failed to get user posts: %v", err)
    }
    defer rows.Close() // **Ferme le curseur après utilisation.

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt) // **Lecture des données du post.
        if err != nil {
            return nil, fmt.Errorf("failed to scan post: %v", err)
        }
        posts = append(posts, post) // **Ajoute le post à la liste des posts.
    }
    return posts, nil // **Retourne la liste des posts.
}

// **AddComment ajoute un commentaire à un post.
func At(db *sql.DB, postID, userID int, content string) error {
    insertCommentSQL := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
    _, err := db.Exec(insertCommentSQL, postID, userID, content) // **Insertion du commentaire dans la BDD.
    if err != nil {
        return fmt.Errorf("failed to add comment: %v", err)
    }
    return nil
}

// **DeletePost supprime un post et ses commentaires associés.
func DeletePost(db *sql.DB, postID, userID int) error {
    tx, err := db.Begin() // **Démarre une transaction.
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }

    // **Supprime les commentaires associés au post.
    _, err = tx.Exec(`DELETE FROM comments WHERE post_id = ?`, postID)
    if err != nil {
        tx.Rollback() // **Annule la transaction en cas d'erreur.
        return fmt.Errorf("failed to delete comments: %v", err)
    }

    // **Delete le post.
    result, err := tx.Exec(`DELETE FROM posts WHERE id = ? AND user_id = ?`, postID, userID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to delete post: %v", err)
    }

    affected, err := result.RowsAffected() // **Vérifie le nombre de lignes affectées.
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to get affected rows: %v", err)
    }

    if affected == 0 {
        tx.Rollback()
        return fmt.Errorf("post not found or user not authorized") // **Aucun post supprimé (inexistant ou mauvaise autorisation).
    }

    return tx.Commit() // **Commit la transaction.
}

// **CloseDB ferme la connexion à la base de données.
func CloseDB(db *sql.DB) {
    if err := db.Close(); err != nil {
        log.Printf("Error closing database: %v", err) // **Erreur lors de la fermeture de la BDD.
    }
}

// **Exporte les données de la base de données vers un fichier texte.
func ExportDataToFile(db *sql.DB, filename string) error {
    file, err := os.Create(filename) // **Crée un fichier pour l'exportation.
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer file.Close() // **Ferme le fichier après utilisation.

    // **Exporte les utilisateurs.
    rows, err := db.Query("SELECT id, username, email, created_at FROM users")
    if err != nil {
        return fmt.Errorf("failed to query users: %v", err)
    }
    defer rows.Close()

    file.WriteString("Users:\n")
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt) // **Lecture des données utilisateur.
        if err != nil {
            return fmt.Errorf("failed to scan user: %v", err)
        }
        file.WriteString(fmt.Sprintf("ID: %d, Username: %s, Email: %s, CreatedAt: %s\n", user.ID, user.Username, user.Email, user.CreatedAt)) // **Écrit les infos utilisateur dans le fichier.
    }

    // **Exporte les posts.
    rows, err = db.Query("SELECT id, user_id, title, content, created_at FROM posts")
    if err != nil {
        return fmt.Errorf("failed to query posts: %v", err)
    }
    defer rows.Close()

    file.WriteString("\nPosts:\n")
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt) // **Lecture des données du post.
        if err != nil {
            return fmt.Errorf("failed to scan post: %v", err)
        }
        file.WriteString(fmt.Sprintf("ID: %d, UserID: %d, Title: %s, Content: %s, CreatedAt: %s\n", post.ID, post.UserID, post.Title, post.Content, post.CreatedAt)) // **Écrit les infos du post dans le fichier.
    }

    rows, err = db.Query("SELECT id, post_id, user_id, content, created_at FROM comments")
    if err != nil {
        return fmt.Errorf("failed to query comments: %v", err)
    }
    defer rows.Close()

    file.WriteString("\nComments:\n")
    for rows.Next() {
        var comment struct {
            ID        int
            PostID    int
            UserID    int
            Content   string
            CreatedAt time.Time
        }
        err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt) // **Lecture des données du commentaire.
        if err != nil {
            return fmt.Errorf("failed to scan comment: %v", err)
        }
        file.WriteString(fmt.Sprintf("ID: %d, PostID: %d, UserID: %d, Content: %s, CreatedAt: %s\n", comment.ID, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)) // **Écrit les infos du commentaire dans le fichier.
    }

    return nil 
}