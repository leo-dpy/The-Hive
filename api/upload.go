package api

import (
	"encoding/base64"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"the-hive/database"
)

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int)

	// Limiter à 5MB
	r.ParseMultipartForm(5 << 20)
	
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		sendError(w, http.StatusBadRequest, "Image invalide ou absente")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		sendError(w, http.StatusBadRequest, "Seules les images (jpg, png, gif) sont autorisées")
		return
	}

	// Déterminer le type MIME
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}
	mimeType := mimeTypes[ext]

	// Lire le fichier en mémoire
	data, err := io.ReadAll(file)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur de lecture de l'image")
		return
	}

	// Convertir en base64 data URI
	b64 := base64.StdEncoding.EncodeToString(data)
	dataURI := "data:" + mimeType + ";base64," + b64

	// Stocker directement dans la base de données
	_, err = database.DB.Exec(`UPDATE users SET profile_picture = ? WHERE id = ?`, dataURI, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur base de données")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{
		"message": "Avatar mis à jour",
		"url":     dataURI,
	})
}
