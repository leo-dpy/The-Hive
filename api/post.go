package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"the-hive/database"
)

const defaultAvatar = `data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxMDAgMTAwIiB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCI+PGNpcmNsZSBjeD0iNTAiIGN5PSI1MCIgcj0iNTAiIGZpbGw9IiMyYTJhMmEiLz48Y2lyY2xlIGN4PSI1MCIgY3k9IjM4IiByPSIxNiIgZmlsbD0iIzU1NSIvPjxlbGxpcHNlIGN4PSI1MCIgY3k9IjgwIiByeD0iMjgiIHJ5PSIyMiIgZmlsbD0iIzU1NSIvPjwvc3ZnPg==`

type CreatePostRequest struct {
	CategoryID *int   `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
}

type PostResponse struct {
	ID             int    `json:"id"`
	Author         string `json:"author"`
	ProfilePicture string `json:"profile_picture"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
	Likes          int    `json:"likes"`
	Dislikes       int    `json:"dislikes"`
	CommentCount   int    `json:"comment_count"`
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Payload invalide")
		return
	}

	if req.Content == "" {
		sendError(w, http.StatusBadRequest, "Le contenu ne peut pas être vide")
		return
	}

	// MySQL: pas de RETURNING id
	res, err := database.DB.Exec(
		`INSERT INTO posts (user_id, category_id, title, content) VALUES (?, ?, ?, ?)`,
		userID, req.CategoryID, req.Title, req.Content,
	)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur lors de la création du post")
		return
	}

	id, _ := res.LastInsertId()

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Post publié !",
		"id":      int(id),
	})
}

func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")

	query := `
		SELECT p.id, u.username, u.profile_picture, p.title, p.content, p.created_at,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = 1) AS likes,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = -1) AS dislikes,
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
`
	var rows *sql.Rows
	var err error

	if categoryID != "" {
		query += ` WHERE p.category_id = ? ORDER BY p.created_at DESC LIMIT 50`
		rows, err = database.DB.Query(query, categoryID)
	} else {
		query += ` ORDER BY p.created_at DESC LIMIT 50`
		rows, err = database.DB.Query(query)
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Impossible de charger les posts")
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
			p.ProfilePicture = defaultAvatar
		}
		posts = append(posts, p)
	}

	if posts == nil {
		posts = []PostResponse{}
	}
	sendJSON(w, http.StatusOK, posts)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		sendError(w, http.StatusBadRequest, "ID du post manquant")
		return
	}

	var p PostResponse
	var title, pic sql.NullString
	
	err := database.DB.QueryRow(`
		SELECT p.id, u.username, u.profile_picture, p.title, p.content, p.created_at,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = 1) AS likes,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'post' AND r.target_id = p.id AND r.value = -1) AS dislikes,
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = ?
	`, postID).Scan(&p.ID, &p.Author, &pic, &title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentCount)

	if err != nil {
		sendError(w, http.StatusNotFound, "Post introuvable")
		return
	}

	p.Title = title.String
	p.ProfilePicture = pic.String
	if p.ProfilePicture == "" {
		p.ProfilePicture = defaultAvatar
	}

	sendJSON(w, http.StatusOK, p)
}
