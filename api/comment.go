package api

import (
	"encoding/json"
	"net/http"
	"the-hive/database"
)

type CreateCommentRequest struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Payload invalide")
		return
	}

	if req.Content == "" {
		sendError(w, http.StatusBadRequest, "Le commentaire ne peut pas être vide")
		return
	}

	_, err := database.DB.Exec(
		`INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`,
		req.PostID, userID, req.Content,
	)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur lors de l'ajout du commentaire")
		return
	}

	sendJSON(w, http.StatusCreated, map[string]string{"message": "Commentaire ajouté !"})
}
