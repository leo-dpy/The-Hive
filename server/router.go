package server

import (
	"log"
	"net/http"

	"the-hive/api"
	"the-hive/config"
)

// Start initializes the HTTP router and starts the server.
func Start() {
	mux := http.NewServeMux()

	// Static files routing
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	// API routes (Publiques)
	mux.HandleFunc("GET /api/ping", api.PingHandler)
	mux.HandleFunc("POST /api/register", api.RegisterHandler)
	mux.HandleFunc("POST /api/login", api.LoginHandler)
	mux.HandleFunc("GET /api/categories", api.GetCategoriesHandler)
	mux.HandleFunc("GET /api/posts", api.ListPostsHandler) // Le Feed
	mux.HandleFunc("GET /api/post", api.GetPostHandler)
	mux.HandleFunc("GET /api/comments", api.GetCommentsHandler)

	// API routes (Protégées par Session)
	mux.HandleFunc("POST /api/logout", api.RequireAuth(api.LogoutHandler))
	mux.HandleFunc("GET /api/me", api.RequireAuth(api.CurrentUserHandler))
	mux.HandleFunc("POST /api/posts", api.RequireAuth(api.CreatePostHandler))
	mux.HandleFunc("POST /api/comments", api.RequireAuth(api.CreateCommentHandler))
	mux.HandleFunc("POST /api/react", api.RequireAuth(api.ReactHandler))
	mux.HandleFunc("POST /api/follow", api.RequireAuth(api.FollowHandler))
	mux.HandleFunc("PUT /api/user/update", api.RequireAuth(api.UpdateUserHandler))
	mux.HandleFunc("POST /api/user/avatar", api.RequireAuth(api.UploadAvatarHandler))
	
	mux.HandleFunc("GET /api/user", api.GetUserProfileHandler)
	mux.HandleFunc("GET /api/user/posts", api.GetUserPostsHandler)
	mux.HandleFunc("GET /api/user/likes", api.GetUserLikesHandler)
	mux.HandleFunc("GET /api/feed/following", api.RequireAuth(api.FeedFollowingHandler))
	mux.HandleFunc("GET /api/messages", api.RequireAuth(api.GetMessagesHandler))
	mux.HandleFunc("GET /api/conversations", api.RequireAuth(api.GetConversationsHandler))
	
	// WebSocket (L'authentification est gérée à l'intérieur du Handler)
	mux.HandleFunc("GET /api/ws", api.WSHandler)

	// Get server address from config
	addr := config.GetEnv("APP_ADDR", ":8080")
	log.Printf("The Hive is running on http://localhost%s", addr)
	
	// Middleware pour logger les requêtes (optionnel mais utile)
	handler := logRequest(mux)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
