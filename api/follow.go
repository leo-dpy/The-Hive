package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"the-hive/database"
)

type FollowRequest struct {
	TargetID int `json:"target_id"`
}

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Payload invalide")
		return
	}

	if req.TargetID == userID {
		sendError(w, http.StatusBadRequest, "Vous ne pouvez pas vous suivre vous-même")
		return
	}

	var count int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ?`,
		userID, req.TargetID,
	).Scan(&count)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur base de données")
		return
	}

	if count > 0 {
		_, err = database.DB.Exec(
			`DELETE FROM followers WHERE follower_id = ? AND following_id = ?`,
			userID, req.TargetID,
		)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Erreur lors du désabonnement")
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "Désabonné avec succès"})
	} else {
		_, err = database.DB.Exec(
			`INSERT INTO followers (follower_id, following_id) VALUES (?, ?)`,
			userID, req.TargetID,
		)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Erreur lors de l'abonnement")
			return
		}
		sendJSON(w, http.StatusCreated, map[string]string{"message": "Abonné avec succès"})
	}
}

func FeedFollowingHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	rows, err := database.DB.Query(`
		SELECT p.id, u.username, u.profile_picture, p.title, p.content, p.created_at,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = 1) AS likes,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = -1) AS dislikes,
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN followers f ON f.following_id = p.user_id
		WHERE f.follower_id = ?
		ORDER BY p.created_at DESC
		LIMIT 50
	`, userID)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Impossible de charger le feed des abonnements")
		return
	}
	defer rows.Close()

	var posts []PostResponse
	for rows.Next() {
		var p PostResponse
		var title, pic sql.NullString
		if err := rows.Scan(&p.ID, &p.Author, &pic, &title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentCount); err != nil {
			continue
		}
		p.Title = title.String
		p.ProfilePicture = pic.String
		if p.ProfilePicture == "" {
			p.ProfilePicture = "/uploads/avatars/default.png"
		}
		posts = append(posts, p)
	}

	if posts == nil {
		posts = []PostResponse{}
	}
	sendJSON(w, http.StatusOK, posts)
}
