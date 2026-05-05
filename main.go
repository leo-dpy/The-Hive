package main

import (
	"the-hive/database"
	"the-hive/server"
)

func main() {
	// 1. Initialisation de la connexion à la base de données (MySQL)
	database.Connect()

	// 2. Insertion des données factices si la base est vide
	database.Seed()

	// 3. Initialisation et démarrage du serveur HTTP
	server.Start()
}
