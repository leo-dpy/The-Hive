package api

import (
	"database/sql"
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

type CommentResponse struct {
	ID             int    `json:"id"`
	Author         string `json:"author"`
	ProfilePicture string `json:"profile_picture"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		sendError(w, http.StatusBadRequest, "post_id manquant")
		return
	}

	rows, err := database.DB.Query(`
		SELECT c.id, u.username, u.profile_picture, c.content, c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur base de données")
		return
	}
	defer rows.Close()

	var comments []CommentResponse
	for rows.Next() {
		var c CommentResponse
		var pic sql.NullString
		if err := rows.Scan(&c.ID, &c.Author, &pic, &c.Content, &c.CreatedAt); err == nil {
			c.ProfilePicture = pic.String
			if c.ProfilePicture == "" { c.ProfilePicture = "" }
			comments = append(comments, c)
		}
	}
	if comments == nil { comments = []CommentResponse{} }
	sendJSON(w, http.StatusOK, comments)
}
