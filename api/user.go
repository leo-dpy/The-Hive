package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"the-hive/database"
)

type UpdateUserRequest struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Bio         string `json:"bio,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Payload invalide")
		return
	}

	var currentUsername, currentEmail, currentHash string
	var currentBio sql.NullString
	err := database.DB.QueryRow(
		`SELECT username, email, bio, password_hash FROM users WHERE id = ?`,
		userID,
	).Scan(&currentUsername, &currentEmail, &currentBio, &currentHash)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur lors de la récupération du profil")
		return
	}

	username := currentUsername
	if req.Username != "" {
		username = strings.TrimSpace(req.Username)
	}

	email := currentEmail
	if req.Email != "" {
		email = strings.TrimSpace(req.Email)
	}

	bio := currentBio.String
	if req.Bio != "" {
		bio = strings.TrimSpace(req.Bio)
	}

	hash := currentHash
	if req.NewPassword != "" {
		if len(req.NewPassword) < 6 {
			sendError(w, http.StatusBadRequest, "Le mot de passe doit faire au moins 6 caractères")
			return
		}
		newHashBytes, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Erreur hash")
			return
		}
		hash = string(newHashBytes)
	}

	_, err = database.DB.Exec(
		`UPDATE users SET username = ?, email = ?, bio = ?, password_hash = ? WHERE id = ?`,
		username, email, bio, hash, userID,
	)

	if err != nil {
		sendError(w, http.StatusConflict, "Pseudo/Email déjà utilisé")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "Profil mis à jour !"})
}

type ProfileResponse struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
	Followers      int    `json:"followers_count"`
	Following      int    `json:"following_count"`
	IsFollowing    bool   `json:"is_following"` // Vue du visiteur
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		sendError(w, http.StatusBadRequest, "Username manquant")
		return
	}

	visitorID := 0
	cookie, err := r.Cookie("hive_session")
	if err == nil {
		database.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, cookie.Value).Scan(&visitorID)
	}

	var p ProfileResponse
	var bio, pic sql.NullString
	err = database.DB.QueryRow(`
		SELECT id, username, bio, profile_picture,
		(SELECT COUNT(*) FROM followers WHERE following_id = users.id) AS followers,
		(SELECT COUNT(*) FROM followers WHERE follower_id = users.id) AS following
		FROM users WHERE username = ?
	`, username).Scan(&p.ID, &p.Username, &bio, &pic, &p.Followers, &p.Following)

	if err != nil {
		sendError(w, http.StatusNotFound, "Profil introuvable")
		return
	}

	p.Bio = bio.String
	p.ProfilePicture = pic.String
	if p.ProfilePicture == "" { p.ProfilePicture = defaultAvatar }

	if visitorID > 0 && visitorID != p.ID {
		var isF int
		database.DB.QueryRow(`SELECT 1 FROM followers WHERE follower_id = ? AND following_id = ?`, visitorID, p.ID).Scan(&isF)
		p.IsFollowing = (isF == 1)
	}

	sendJSON(w, http.StatusOK, p)
}

func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	rows, err := database.DB.Query(`
		SELECT p.id, u.username, u.profile_picture, p.title, p.content, p.created_at,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = 1) AS likes,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = -1) AS dislikes,
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE u.username = ?
		ORDER BY p.created_at DESC
		LIMIT 50
	`, username)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur DB")
		return
	}
	defer rows.Close()

	var posts []PostResponse
	for rows.Next() {
		var p PostResponse
		var title, pic sql.NullString
		if err := rows.Scan(&p.ID, &p.Author, &pic, &title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentCount); err == nil {
			p.Title = title.String
			p.ProfilePicture = pic.String
			if p.ProfilePicture == "" { p.ProfilePicture = defaultAvatar }
			posts = append(posts, p)
		}
	}
	if posts == nil { posts = []PostResponse{} }
	sendJSON(w, http.StatusOK, posts)
}

func GetUserLikesHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	rows, err := database.DB.Query(`
		SELECT p.id, u.username, u.profile_picture, p.title, p.content, p.created_at,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = 1) AS likes,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = -1) AS dislikes,
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN reactions r ON r.target_id = p.id AND r.target_type = 'post' AND r.value = 1
		JOIN users liker ON liker.id = r.user_id
		WHERE liker.username = ?
		ORDER BY r.created_at DESC
		LIMIT 50
	`, username)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur DB")
		return
	}
	defer rows.Close()

	var posts []PostResponse
	for rows.Next() {
		var p PostResponse
		var title, pic sql.NullString
		if err := rows.Scan(&p.ID, &p.Author, &pic, &title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentCount); err == nil {
			p.Title = title.String
			p.ProfilePicture = pic.String
			if p.ProfilePicture == "" { p.ProfilePicture = "" }
			posts = append(posts, p)
		}
	}
	if posts == nil { posts = []PostResponse{} }
	sendJSON(w, http.StatusOK, posts)
}
