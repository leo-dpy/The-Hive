package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

const (
	sessionCookieName = "hive_session"
	sessionDuration   = 7 * 24 * time.Hour
)

type App struct {
	db        *sql.DB
	templates *template.Template
}

type User struct {
	ID       int
	Username string
	Email    string
}

type Category struct {
	ID          int
	Name        string
	Description string
}

type ThreadView struct {
	ID           int
	CategoryID   int
	CategoryName string
	Title        string
	Content      string
	CreatedAt    string
	Author       string
	Likes        int
	Dislikes     int
	CommentCount int
	UserReaction int
}

type CommentView struct {
	ID           int
	ThreadID     int
	Content      string
	CreatedAt    string
	Author       string
	Likes        int
	Dislikes     int
	UserReaction int
}

type PageData struct {
	Title           string
	User            *User
	CurrentPath     string
	Flash           string
	Error           string
	Categories      []Category
	Category        *Category
	Threads         []ThreadView
	Thread          *ThreadView
	Comments        []CommentView
	DashboardFilter string
}

func main() {
	if err := os.MkdirAll("data", 0o755); err != nil {
		log.Fatalf("création du dossier data: %v", err)
	}

	dbPath := filepath.ToSlash(filepath.Join("data", "hive.db"))
	db, err := sql.Open("sqlite", "file:"+dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("ouverture sqlite: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	app := &App{db: db}
	if err := app.initDB(); err != nil {
		log.Fatalf("initialisation base de données: %v", err)
	}
	if err := app.loadTemplates(); err != nil {
		log.Fatalf("chargement templates: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/category", app.categoryHandler)
	mux.HandleFunc("/thread", app.threadHandler)
	mux.HandleFunc("/thread/create", app.createThreadHandler)
	mux.HandleFunc("/dashboard", app.dashboardHandler)
	mux.HandleFunc("/login", app.loginPageHandler)
	mux.HandleFunc("/register", app.registerPageHandler)

	mux.HandleFunc("/auth/register", app.registerActionHandler)
	mux.HandleFunc("/auth/login", app.loginActionHandler)
	mux.HandleFunc("/auth/logout", app.logoutActionHandler)
	mux.HandleFunc("/comment/add", app.addCommentHandler)
	mux.HandleFunc("/react", app.reactHandler)

	addr := getEnv("APP_ADDR", ":8080")
	log.Printf("The Hive est lancé sur http://localhost%s", addr)
	if err := http.ListenAndServe(addr, app.recoverMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func (a *App) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *App) initDB() error {
	statements := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		);`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS threads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			thread_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			target_type TEXT NOT NULL CHECK (target_type IN ('thread', 'comment')),
			target_id INTEGER NOT NULL,
			value INTEGER NOT NULL CHECK (value IN (1, -1)),
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE (user_id, target_type, target_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_threads_category ON threads(category_id);`,
		`CREATE INDEX IF NOT EXISTS idx_threads_user ON threads(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_thread ON comments(thread_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reactions_target ON reactions(target_type, target_id);`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);`,
	}

	for _, statement := range statements {
		if _, err := a.db.Exec(statement); err != nil {
			return err
		}
	}

	return a.seedCategories()
}

func (a *App) seedCategories() error {
	var count int
	if err := a.db.QueryRow(`SELECT COUNT(*) FROM categories`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	categories := []Category{
		{Name: "Général", Description: "Discussions libres autour de The Hive."},
		{Name: "Golang", Description: "Questions et ressources sur le backend Go."},
		{Name: "SQLite", Description: "Schéma, requêtes et administration de la base."},
		{Name: "Frontend", Description: "HTML, CSS et JavaScript Vanilla."},
		{Name: "Annonces", Description: "Nouveautés importantes de la communauté."},
	}
	for _, category := range categories {
		if _, err := a.db.Exec(`INSERT INTO categories (name, description) VALUES (?, ?)`, category.Name, category.Description); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) loadTemplates() error {
	funcs := template.FuncMap{
		"reactionClass": func(current int, value int) string {
			if current == value {
				return "is-active"
			}
			return ""
		},
	}

	t, err := template.New("the-hive").Funcs(funcs).ParseGlob(filepath.Join("templates", "*.tmpl"))
	if err != nil {
		return err
	}
	a.templates = t
	return nil
}

func (a *App) render(w http.ResponseWriter, r *http.Request, name string, data PageData) {
	user, _ := a.currentUser(r)
	data.User = user
	data.CurrentPath = r.URL.RequestURI()
	data.Flash = flashMessage(r.URL.Query().Get("ok"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
		http.Error(w, "Erreur de rendu", http.StatusInternalServerError)
	}
}

func flashMessage(code string) string {
	switch code {
	case "registered":
		return "Compte créé avec succès. Bienvenue dans The Hive !"
	case "login":
		return "Connexion réussie."
	case "logout":
		return "Déconnexion réussie."
	case "thread-created":
		return "Sujet publié avec succès."
	case "comment-added":
		return "Commentaire ajouté."
	case "reaction":
		return "Réaction enregistrée."
	default:
		return ""
	}
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	userID := currentUserID(a, r)
	categories, err := a.getCategories()
	if err != nil {
		http.Error(w, "Impossible de charger les catégories", http.StatusInternalServerError)
		return
	}
	threads, err := a.listThreads(userID, "", "LIMIT 20")
	if err != nil {
		http.Error(w, "Impossible de charger les sujets", http.StatusInternalServerError)
		return
	}
	a.render(w, r, "home", PageData{Title: "Accueil", Categories: categories, Threads: threads})
}

func (a *App) categoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	id, err := positiveID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Catégorie invalide", http.StatusBadRequest)
		return
	}
	category, err := a.getCategory(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Impossible de charger la catégorie", http.StatusInternalServerError)
		return
	}
	userID := currentUserID(a, r)
	threads, err := a.listThreads(userID, "t.category_id = ?", "", id)
	if err != nil {
		http.Error(w, "Impossible de charger les sujets", http.StatusInternalServerError)
		return
	}
	a.render(w, r, "category", PageData{Title: category.Name, Category: category, Threads: threads})
}

func (a *App) threadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	id, err := positiveID(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Sujet invalide", http.StatusBadRequest)
		return
	}
	userID := currentUserID(a, r)
	thread, err := a.getThread(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Impossible de charger le sujet", http.StatusInternalServerError)
		return
	}
	comments, err := a.getComments(id, userID)
	if err != nil {
		http.Error(w, "Impossible de charger les commentaires", http.StatusInternalServerError)
		return
	}
	a.render(w, r, "thread", PageData{Title: thread.Title, Thread: thread, Comments: comments})
}

func (a *App) createThreadHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	categories, err := a.getCategories()
	if err != nil {
		http.Error(w, "Impossible de charger les catégories", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.render(w, r, "create_thread", PageData{Title: "Créer un sujet", Categories: categories})
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulaire invalide", http.StatusBadRequest)
			return
		}
		categoryID, err := positiveID(r.FormValue("category_id"))
		if err != nil {
			a.render(w, r, "create_thread", PageData{Title: "Créer un sujet", Categories: categories, Error: "Choisis une catégorie valide."})
			return
		}
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		if len(title) < 3 || len(title) > 140 {
			a.render(w, r, "create_thread", PageData{Title: "Créer un sujet", Categories: categories, Error: "Le titre doit contenir entre 3 et 140 caractères."})
			return
		}
		if len(content) < 10 {
			a.render(w, r, "create_thread", PageData{Title: "Créer un sujet", Categories: categories, Error: "Le contenu doit contenir au moins 10 caractères."})
			return
		}
		res, err := a.db.Exec(`INSERT INTO threads (user_id, category_id, title, content) VALUES (?, ?, ?, ?)`, user.ID, categoryID, title, content)
		if err != nil {
			log.Printf("insert thread: %v", err)
			a.render(w, r, "create_thread", PageData{Title: "Créer un sujet", Categories: categories, Error: "Impossible de créer le sujet. Vérifie la catégorie."})
			return
		}
		threadID, _ := res.LastInsertId()
		http.Redirect(w, r, fmt.Sprintf("/thread?id=%d&ok=thread-created", threadID), http.StatusSeeOther)
	default:
		methodNotAllowed(w)
	}
}

func (a *App) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "posted"
	}

	var (
		threads []ThreadView
		err     error
	)
	switch filter {
	case "posted":
		threads, err = a.listThreads(user.ID, "t.user_id = ?", "", user.ID)
	case "liked":
		threads, err = a.listThreads(user.ID, `EXISTS (
			SELECT 1 FROM reactions rr
			WHERE rr.target_type = 'thread'
			AND rr.target_id = t.id
			AND rr.user_id = ?
			AND rr.value = 1
		)`, "", user.ID)
	default:
		http.Error(w, "Filtre inconnu", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Impossible de charger le dashboard", http.StatusInternalServerError)
		return
	}
	a.render(w, r, "dashboard", PageData{Title: "Dashboard", Threads: threads, DashboardFilter: filter})
}

func (a *App) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	if _, ok := a.currentUser(r); ok {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	a.render(w, r, "login", PageData{Title: "Connexion"})
}

func (a *App) registerPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	if _, ok := a.currentUser(r); ok {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	a.render(w, r, "register", PageData{Title: "Inscription"})
}

func (a *App) registerActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	password := r.FormValue("password")

	if len(username) < 3 || len(username) > 32 || strings.ContainsAny(username, " <>/\\") {
		a.render(w, r, "register", PageData{Title: "Inscription", Error: "Le pseudo doit contenir 3 à 32 caractères sans espace ni symbole HTML."})
		return
	}
	if !strings.Contains(email, "@") || len(email) < 6 {
		a.render(w, r, "register", PageData{Title: "Inscription", Error: "Adresse email invalide."})
		return
	}
	if len(password) < 8 {
		a.render(w, r, "register", PageData{Title: "Inscription", Error: "Le mot de passe doit contenir au moins 8 caractères."})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Impossible de sécuriser le mot de passe", http.StatusInternalServerError)
		return
	}
	res, err := a.db.Exec(`INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`, username, email, string(hash))
	if err != nil {
		log.Printf("insert user: %v", err)
		a.render(w, r, "register", PageData{Title: "Inscription", Error: "Pseudo ou email déjà utilisé."})
		return
	}
	userID, _ := res.LastInsertId()
	if err := a.createSession(w, int(userID)); err != nil {
		http.Error(w, "Compte créé, mais connexion impossible", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard?ok=registered", http.StatusSeeOther)
}

func (a *App) loginActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}
	identifier := strings.TrimSpace(r.FormValue("identifier"))
	password := r.FormValue("password")
	var (
		userID       int
		passwordHash string
	)
	err := a.db.QueryRow(`SELECT id, password_hash FROM users WHERE username = ? OR email = ?`, identifier, strings.ToLower(identifier)).Scan(&userID, &passwordHash)
	if err != nil {
		a.render(w, r, "login", PageData{Title: "Connexion", Error: "Identifiants invalides."})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		a.render(w, r, "login", PageData{Title: "Connexion", Error: "Identifiants invalides."})
		return
	}
	if err := a.createSession(w, userID); err != nil {
		http.Error(w, "Impossible de créer la session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard?ok=login", http.StatusSeeOther)
}

func (a *App) logoutActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		_, _ = a.db.Exec(`DELETE FROM sessions WHERE id = ?`, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/?ok=logout", http.StatusSeeOther)
}

func (a *App) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}
	threadID, err := positiveID(r.FormValue("thread_id"))
	if err != nil {
		http.Error(w, "Sujet invalide", http.StatusBadRequest)
		return
	}
	content := strings.TrimSpace(r.FormValue("content"))
	if len(content) < 2 {
		http.Redirect(w, r, fmt.Sprintf("/thread?id=%d", threadID), http.StatusSeeOther)
		return
	}
	if _, err := a.db.Exec(`INSERT INTO comments (thread_id, user_id, content) VALUES (?, ?, ?)`, threadID, user.ID, content); err != nil {
		log.Printf("insert comment: %v", err)
		http.Error(w, "Impossible d'ajouter le commentaire", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/thread?id=%d&ok=comment-added", threadID), http.StatusSeeOther)
}

func (a *App) reactHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}
	targetType := r.FormValue("target_type")
	if targetType != "thread" && targetType != "comment" {
		http.Error(w, "Type de cible invalide", http.StatusBadRequest)
		return
	}
	targetID, err := positiveID(r.FormValue("target_id"))
	if err != nil {
		http.Error(w, "Cible invalide", http.StatusBadRequest)
		return
	}
	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil || (value != 1 && value != -1) {
		http.Error(w, "Réaction invalide", http.StatusBadRequest)
		return
	}
	if exists, err := a.targetExists(targetType, targetID); err != nil || !exists {
		if err != nil {
			log.Printf("target exists: %v", err)
		}
		http.Error(w, "Cible inexistante", http.StatusBadRequest)
		return
	}

	var oldValue int
	err = a.db.QueryRow(`SELECT value FROM reactions WHERE user_id = ? AND target_type = ? AND target_id = ?`, user.ID, targetType, targetID).Scan(&oldValue)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = a.db.Exec(`INSERT INTO reactions (user_id, target_type, target_id, value) VALUES (?, ?, ?, ?)`, user.ID, targetType, targetID, value)
	case err == nil && oldValue == value:
		_, err = a.db.Exec(`DELETE FROM reactions WHERE user_id = ? AND target_type = ? AND target_id = ?`, user.ID, targetType, targetID)
	case err == nil:
		_, err = a.db.Exec(`UPDATE reactions SET value = ?, created_at = datetime('now') WHERE user_id = ? AND target_type = ? AND target_id = ?`, value, user.ID, targetType, targetID)
	}
	if err != nil {
		log.Printf("react: %v", err)
		http.Error(w, "Impossible d'enregistrer la réaction", http.StatusInternalServerError)
		return
	}
	redirect := safeRedirect(r.FormValue("redirect"))
	if !strings.Contains(redirect, "ok=") {
		separator := "?"
		if strings.Contains(redirect, "?") {
			separator = "&"
		}
		redirect += separator + "ok=reaction"
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (a *App) getCategories() ([]Category, error) {
	rows, err := a.db.Query(`SELECT id, name, description FROM categories ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (a *App) getCategory(id int) (*Category, error) {
	var category Category
	err := a.db.QueryRow(`SELECT id, name, description FROM categories WHERE id = ?`, id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (a *App) listThreads(userID int, where string, suffix string, args ...any) ([]ThreadView, error) {
	query := `SELECT
		t.id,
		t.category_id,
		c.name,
		t.title,
		t.content,
		t.created_at,
		u.username,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'thread' AND r.target_id = t.id AND r.value = 1) AS likes,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'thread' AND r.target_id = t.id AND r.value = -1) AS dislikes,
		(SELECT COUNT(*) FROM comments cm WHERE cm.thread_id = t.id) AS comment_count,
		COALESCE((SELECT r.value FROM reactions r WHERE r.target_type = 'thread' AND r.target_id = t.id AND r.user_id = ?), 0) AS user_reaction
	FROM threads t
	JOIN users u ON u.id = t.user_id
	JOIN categories c ON c.id = t.category_id`
	queryArgs := []any{userID}
	if where != "" {
		query += " WHERE " + where
		queryArgs = append(queryArgs, args...)
	}
	query += " ORDER BY t.created_at DESC, t.id DESC " + suffix

	rows, err := a.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []ThreadView
	for rows.Next() {
		var thread ThreadView
		if err := rows.Scan(
			&thread.ID,
			&thread.CategoryID,
			&thread.CategoryName,
			&thread.Title,
			&thread.Content,
			&thread.CreatedAt,
			&thread.Author,
			&thread.Likes,
			&thread.Dislikes,
			&thread.CommentCount,
			&thread.UserReaction,
		); err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	return threads, rows.Err()
}

func (a *App) getThread(id int, userID int) (*ThreadView, error) {
	threads, err := a.listThreads(userID, "t.id = ?", "", id)
	if err != nil {
		return nil, err
	}
	if len(threads) == 0 {
		return nil, sql.ErrNoRows
	}
	return &threads[0], nil
}

func (a *App) getComments(threadID int, userID int) ([]CommentView, error) {
	rows, err := a.db.Query(`SELECT
		cm.id,
		cm.thread_id,
		cm.content,
		cm.created_at,
		u.username,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'comment' AND r.target_id = cm.id AND r.value = 1) AS likes,
		(SELECT COUNT(*) FROM reactions r WHERE r.target_type = 'comment' AND r.target_id = cm.id AND r.value = -1) AS dislikes,
		COALESCE((SELECT r.value FROM reactions r WHERE r.target_type = 'comment' AND r.target_id = cm.id AND r.user_id = ?), 0) AS user_reaction
	FROM comments cm
	JOIN users u ON u.id = cm.user_id
	WHERE cm.thread_id = ?
	ORDER BY cm.created_at ASC, cm.id ASC`, userID, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentView
	for rows.Next() {
		var comment CommentView
		if err := rows.Scan(
			&comment.ID,
			&comment.ThreadID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.Author,
			&comment.Likes,
			&comment.Dislikes,
			&comment.UserReaction,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (a *App) createSession(w http.ResponseWriter, userID int) error {
	sessionID := uuid.NewString()
	expiresAt := time.Now().UTC().Add(sessionDuration)
	if _, err := a.db.Exec(`DELETE FROM sessions WHERE expires_at <= ?`, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return err
	}
	if _, err := a.db.Exec(`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`, sessionID, userID, expiresAt.Format(time.RFC3339)); err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(sessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func (a *App) currentUser(r *http.Request) (*User, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, false
	}
	var (
		user      User
		expiresAt string
	)
	err = a.db.QueryRow(`SELECT u.id, u.username, u.email, s.expires_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = ?`, cookie.Value).Scan(&user.ID, &user.Username, &user.Email, &expiresAt)
	if err != nil {
		return nil, false
	}
	expiry, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil || time.Now().UTC().After(expiry) {
		_, _ = a.db.Exec(`DELETE FROM sessions WHERE id = ?`, cookie.Value)
		return nil, false
	}
	return &user, true
}

func currentUserID(a *App, r *http.Request) int {
	user, ok := a.currentUser(r)
	if !ok {
		return 0
	}
	return user.ID
}

func (a *App) targetExists(targetType string, targetID int) (bool, error) {
	var count int
	var err error
	if targetType == "thread" {
		err = a.db.QueryRow(`SELECT COUNT(*) FROM threads WHERE id = ?`, targetID).Scan(&count)
	} else {
		err = a.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE id = ?`, targetID).Scan(&count)
	}
	return count > 0, err
}

func positiveID(raw string) (int, error) {
	id, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}

func safeRedirect(raw string) string {
	if raw == "" || !strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "//") {
		return "/"
	}
	return raw
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Allow", "GET, POST")
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}
