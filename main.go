package main

import (
	"the-hive/database"
	"the-hive/server"
)

func main() {
	// 1. Initialisation de la connexion à la base de données (PostgreSQL)
	database.Connect()

	// 2. Initialisation et démarrage du serveur HTTP
	server.Start()
}
