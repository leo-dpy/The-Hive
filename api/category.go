package api

import (
	"net/http"
	"the-hive/database"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`SELECT id, name, description FROM categories ORDER BY id ASC`)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur lors du chargement des catégories")
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err == nil {
			categories = append(categories, c)
		}
	}
	if categories == nil {
		categories = []Category{}
	}
	sendJSON(w, http.StatusOK, categories)
}