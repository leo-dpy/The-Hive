package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"
	"the-hive/database"
)

type contextKey string
const UserIDKey contextKey = "userID"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("hive_session")
		if err != nil {
			sendError(w, http.StatusUnauthorized, "Non autorisé")
			return
		}

		var userID int
		var expiresAt time.Time
		// MySQL: ? au lieu de $1
		err = database.DB.QueryRow(
			`SELECT user_id, expires_at FROM sessions WHERE id = ?`,
			cookie.Value,
		).Scan(&userID, &expiresAt)

		if err != nil || time.Now().After(expiresAt) {
			sendError(w, http.StatusUnauthorized, "Session invalide ou expirée")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	var username, email string
	var bio, profilePicture sql.NullString
	err := database.DB.QueryRow(`SELECT username, email, bio, profile_picture FROM users WHERE id = ?`, userID).Scan(&username, &email, &bio, &profilePicture)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Utilisateur introuvable")
		return
	}

	pic := profilePicture.String
	if pic == "" {
		pic = defaultAvatar
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"id":              userID,
		"username":        username,
		"email":           email,
		"bio":             bio.String,
		"profile_picture": pic,
	})
}
