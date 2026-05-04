package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// init est appelé automatiquement au démarrage du package.
// Il va tenter de charger le fichier .env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Aucun fichier .env trouvé ou erreur de lecture, utilisation des variables système.")
	} else {
		log.Println("✅ Fichier .env chargé avec succès.")
	}
}

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the fallback value if the variable is not present or is empty.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}
