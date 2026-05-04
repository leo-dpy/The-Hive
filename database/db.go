package database

import (
	"database/sql"
	"log"
	"net/url"
	"os"
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

func initSchema() {
	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Printf("⚠️ Fichier database/schema.sql introuvable.")
		return
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatalf("❌ Erreur lors de l'exécution du schéma MySQL: %v", err)
	}
	
	log.Println("✅ Schéma MySQL initialisé avec succès.")
}
