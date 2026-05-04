package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"the-hive/database"
)

type ReactionRequest struct {
	TargetType string `json:"target_type"` // "post" ou "comment"
	TargetID   int    `json:"target_id"`
	Value      int    `json:"value"` // 1 pour Upvote, -1 pour Downvote
}

func ReactHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	var req ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Payload invalide")
		return
	}

	if req.TargetType != "post" && req.TargetType != "comment" {
		sendError(w, http.StatusBadRequest, "Target type invalide")
		return
	}

	if req.Value != 1 && req.Value != -1 {
		sendError(w, http.StatusBadRequest, "Value invalide (1 ou -1)")
		return
	}

	var existingValue int
	err := database.DB.QueryRow(
		`SELECT value FROM reactions WHERE user_id = ? AND target_type = ? AND target_id = ?`,
		userID, req.TargetType, req.TargetID,
	).Scan(&existingValue)

	if err == sql.ErrNoRows {
		_, err = database.DB.Exec(
			`INSERT INTO reactions (user_id, target_type, target_id, value) VALUES (?, ?, ?, ?)`,
			userID, req.TargetType, req.TargetID, req.Value,
		)
	} else if err == nil {
		if existingValue == req.Value {
			_, err = database.DB.Exec(
				`DELETE FROM reactions WHERE user_id = ? AND target_type = ? AND target_id = ?`,
				userID, req.TargetType, req.TargetID,
			)
		} else {
			_, err = database.DB.Exec(
				`UPDATE reactions SET value = ? WHERE user_id = ? AND target_type = ? AND target_id = ?`,
				req.Value, userID, req.TargetType, req.TargetID,
			)
		}
	}

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur lors de la réaction")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "Réaction mise à jour !"})
}
