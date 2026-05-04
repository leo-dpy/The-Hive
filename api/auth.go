package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"the-hive/database"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"` // Peut être email ou username
	Password string `json:"password"`
}

func sendJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	if len(req.Username) < 3 || len(req.Password) < 6 || req.Email == "" {
		sendError(w, http.StatusBadRequest, "Données invalides")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur serveur (hash)")
		return
	}

	// MySQL: Utilisation de ? au lieu de $1, et LastInsertId() car pas de RETURNING
	res, err := database.DB.Exec(
		`INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`,
		req.Username, req.Email, string(hash),
	)

	if err != nil {
		sendError(w, http.StatusConflict, "Ce nom d'utilisateur ou cet email est déjà pris")
		return
	}

	id, _ := res.LastInsertId()
	userID := int(id)

	createSession(w, userID)
	sendJSON(w, http.StatusCreated, map[string]string{"message": "Inscription réussie !"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	var userID int
	var hash string
	// MySQL: Utilisation de ?
	err := database.DB.QueryRow(
		`SELECT id, password_hash FROM users WHERE username = ? OR email = ?`,
		req.Username, req.Username,
	).Scan(&userID, &hash)

	if err != nil {
		sendError(w, http.StatusUnauthorized, "Identifiants invalides")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		sendError(w, http.StatusUnauthorized, "Identifiants invalides")
		return
	}

	createSession(w, userID)
	sendJSON(w, http.StatusOK, map[string]string{"message": "Connexion réussie !"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("hive_session")
	if err == nil {
		database.DB.Exec(`DELETE FROM sessions WHERE id = ?`, cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "hive_session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	sendJSON(w, http.StatusOK, map[string]string{"message": "Déconnexion réussie !"})
}

func createSession(w http.ResponseWriter, userID int) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err := database.DB.Exec(
		`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`,
		sessionID, userID, expiresAt,
	)
	if err != nil {
		log.Printf("Erreur création de session DB : %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "hive_session",
		Value:    sessionID,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
