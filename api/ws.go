package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"the-hive/database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepte toutes les requêtes pour le développement
	},
}

// Hub stocke les connexions actives par ID utilisateur
type Hub struct {
	sync.RWMutex
	clients map[int]*websocket.Conn
}

var GlobalHub = &Hub{
	clients: make(map[int]*websocket.Conn),
}

// WSMessage est le payload JSON reçu et envoyé via WebSocket
type WSMessage struct {
	ReceiverID int    `json:"receiver_id,omitempty"`
	SenderID   int    `json:"sender_id,omitempty"`
	Content    string `json:"content"`
	Type       string `json:"type,omitempty"`
}

// WSHandler gère l'Upgrade HTTP vers WebSocket
func WSHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Authentification manuelle (pas de middleware car c'est un upgrade)
	cookie, err := r.Cookie("hive_session")
	if err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	var userID int
	err = database.DB.QueryRow(
		`SELECT user_id FROM sessions WHERE id = ?`,
		cookie.Value,
	).Scan(&userID)

	if err != nil {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	// 2. Upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erreur WebSocket Upgrade:", err)
		return
	}

	// 3. Ajouter au Hub
	GlobalHub.Lock()
	GlobalHub.clients[userID] = conn
	GlobalHub.Unlock()

	defer func() {
		GlobalHub.Lock()
		delete(GlobalHub.clients, userID)
		GlobalHub.Unlock()
		conn.Close()
	}()

	// 4. Écoute des messages entrants
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			break // Connexion fermée
		}

		var msg WSMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			continue
		}

		// Sauvegarde en DB MySQL
		_, err = database.DB.Exec(
			`INSERT INTO direct_messages (sender_id, receiver_id, content) VALUES (?, ?, ?)`,
			userID, msg.ReceiverID, msg.Content,
		)
		if err != nil {
			log.Println("Erreur lors de la sauvegarde du DM:", err)
			continue
		}

		// Push instantané au destinataire s'il est en ligne
		GlobalHub.RLock()
		receiverConn, isOnline := GlobalHub.clients[msg.ReceiverID]
		GlobalHub.RUnlock()

		if isOnline {
			response := WSMessage{
				Type:     "new_dm",
				SenderID: userID,
				Content:  msg.Content,
			}
			responseJSON, _ := json.Marshal(response)
			receiverConn.WriteMessage(websocket.TextMessage, responseJSON)
		}
	}
}

// GetMessagesHandler renvoie l'historique de chat
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)
	targetID := r.URL.Query().Get("user_id")

	if targetID == "" {
		sendError(w, http.StatusBadRequest, "Missing user_id parameter")
		return
	}

	rows, err := database.DB.Query(`
		SELECT sender_id, receiver_id, content, created_at
		FROM direct_messages
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC
	`, userID, targetID, targetID, userID)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur base de données")
		return
	}
	defer rows.Close()

	type DM struct {
		SenderID   int    `json:"sender_id"`
		ReceiverID int    `json:"receiver_id"`
		Content    string `json:"content"`
		CreatedAt  string `json:"created_at"`
	}

	var dms []DM
	for rows.Next() {
		var dm DM
		if err := rows.Scan(&dm.SenderID, &dm.ReceiverID, &dm.Content, &dm.CreatedAt); err == nil {
			dms = append(dms, dm)
		}
	}

	if dms == nil {
		dms = []DM{}
	}
	sendJSON(w, http.StatusOK, dms)
}

// GetConversationsHandler renvoie la liste des contacts récents
func GetConversationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	rows, err := database.DB.Query(`
		SELECT u.id, u.username, u.profile_picture,
		  (SELECT content FROM direct_messages 
		   WHERE (sender_id = u.id AND receiver_id = ?) OR (receiver_id = u.id AND sender_id = ?) 
		   ORDER BY created_at DESC LIMIT 1) as last_msg,
		  (SELECT created_at FROM direct_messages 
		   WHERE (sender_id = u.id AND receiver_id = ?) OR (receiver_id = u.id AND sender_id = ?) 
		   ORDER BY created_at DESC LIMIT 1) as last_time
		FROM users u
		WHERE EXISTS (
		   SELECT 1 FROM direct_messages 
		   WHERE (sender_id = u.id AND receiver_id = ?) OR (receiver_id = u.id AND sender_id = ?)
		) AND u.id != ?
		ORDER BY last_time DESC
	`, userID, userID, userID, userID, userID, userID, userID)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur base de données")
		return
	}
	defer rows.Close()

	type Conversation struct {
		UserID         int    `json:"user_id"`
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
		LastMessage    string `json:"last_message"`
		LastTime       string `json:"last_time"`
	}

	var convos []Conversation
	for rows.Next() {
		var c Conversation
		var pic, lastMsg, lastTime sql.NullString
		if err := rows.Scan(&c.UserID, &c.Username, &pic, &lastMsg, &lastTime); err == nil {
			c.ProfilePicture = pic.String
			if c.ProfilePicture == "" { c.ProfilePicture = "" }
			c.LastMessage = lastMsg.String
			c.LastTime = lastTime.String
			convos = append(convos, c)
		}
	}

	if convos == nil {
		convos = []Conversation{}
	}
	sendJSON(w, http.StatusOK, convos)
}
