package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"the-hive/database"
	"time"
	"fmt"
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

	ext := filepath.Ext(handler.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		sendError(w, http.StatusBadRequest, "Seules les images (jpg, png, gif) sont autorisées")
		return
	}

	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	savePath := filepath.Join("public", "uploads", "avatars", filename)
	dbPath := "/uploads/avatars/" + filename

	dst, err := os.Create(savePath)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur lors de la sauvegarde")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur d'écriture de l'image")
		return
	}

	_, err = database.DB.Exec(`UPDATE users SET profile_picture = ? WHERE id = ?`, dbPath, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Erreur base de données")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{
		"message": "Avatar mis à jour",
		"url":     dbPath,
	})
}
