package database

import (
	"database/sql"
	"log"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"the-hive/config"
)

// DB est l'instance globale de la connexion à la base de données.
var DB *sql.DB

// Connect initialise la connexion à la base de données MySQL.
func Connect() {
	dbURL := config.GetEnv("DATABASE_URL", "mysql://root:root@127.0.0.1:3306/hive_db")

	dsn := dbURL
	// Conversion de l'URL mysql:// au format DSN requis par le driver go-sql-driver/mysql
	if strings.HasPrefix(dbURL, "mysql://") {
		u, err := url.Parse(dbURL)
		if err == nil {
			password, _ := u.User.Password()
			// user:password@tcp(host:port)/dbname?parseTime=true&multiStatements=true
			dsn = u.User.Username() + ":" + password + "@tcp(" + u.Host + ")" + u.Path + "?parseTime=true&multiStatements=true"
		}
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Erreur lors de l'ouverture de la connexion MySQL: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Impossible de joindre la base MySQL: %v", err)
	}

	log.Println("✅ Connexion à MySQL réussie !")
	initSchema()
}

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    profile_picture VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id CHAR(36) PRIMARY KEY,
    user_id INT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    category_id INT,
    title VARCHAR(150),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_id INT NOT NULL,
    value INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, target_type, target_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (target_type IN ('post', 'comment')),
    CHECK (value IN (1, -1))
);

CREATE TABLE IF NOT EXISTS direct_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS followers (
    follower_id INT NOT NULL,
    following_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (follower_id != following_id)
);
`

func initSchema() {
	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("❌ Erreur lors de l'exécution du schéma MySQL: %v", err)
	}
	log.Println("✅ Schéma MySQL initialisé avec succès.")
}

// Seed is a placeholder function to avoid compilation errors.
func Seed() {
	log.Println("🌱 Seeding skipped/stubbed.")
}

