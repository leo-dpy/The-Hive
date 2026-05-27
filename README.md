# 🐝 The Hive — Réseau Social

Réseau social moderne développé en **Go** avec une API REST, une messagerie en temps réel par **WebSocket**, et un frontend **HTML/CSS/JS Vanilla** au design glassmorphism.

## Fonctionnalités

### Public
- Consultation du feed global (posts récents)
- Exploration des catégories et des sujets
- Consultation des profils utilisateurs

### Membres
- Inscription / Connexion sécurisée (bcrypt + sessions UUID)
- Publication de posts avec sélection de catégorie
- Commentaires sur les posts
- Système de votes (upvote / downvote) sur les posts
- Suivi d'utilisateurs (follow / unfollow)
- Feed personnalisé basé sur les abonnements (avec fallback sur le feed public)
- Messagerie privée en temps réel (WebSocket)
- Modification du profil (pseudo, email, bio, avatar, mot de passe)

## Stack technique

| Composant | Technologie |
|---|---|
| Backend | Go (net/http, standard library) |
| Base de données | MySQL |
| Temps réel | WebSocket (gorilla/websocket) |
| Auth | bcrypt + sessions UUID cookie |
| Frontend | HTML5, CSS3 (glassmorphism), JavaScript Vanilla |
| Déploiement | Docker multi-stage, Coolify |

### Dépendances Go

```
github.com/go-sql-driver/mysql
github.com/google/uuid
github.com/gorilla/websocket
github.com/joho/godotenv
golang.org/x/crypto/bcrypt
```

## Architecture

```
the-hive/
├── main.go              # Point d'entrée
├── server/
│   └── router.go        # Routeur HTTP et middleware
├── api/
│   ├── auth.go          # Inscription, connexion, déconnexion
│   ├── middleware.go     # Middleware d'authentification
│   ├── post.go          # CRUD posts
│   ├── comment.go       # CRUD commentaires
│   ├── reaction.go      # Likes / dislikes
│   ├── follow.go        # Follow / unfollow + feed abonnements
│   ├── user.go          # Profil utilisateur
│   ├── upload.go        # Upload d'avatars
│   ├── category.go      # Catégories
│   ├── ws.go            # WebSocket messagerie
│   └── ping.go          # Health check
├── config/
│   └── env.go           # Configuration .env
├── database/
│   └── db.go            # Connexion MySQL + migrations
├── public/
│   ├── css/style.css    # Design glassmorphism
│   ├── js/
│   │   ├── auth-modal.js    # Modale de connexion/inscription
│   │   ├── sidebar.js       # Navigation sidebar
│   │   ├── home.js          # Feed principal
│   │   ├── explore.js       # Exploration catégories/posts
│   │   ├── profile.js       # Page profil
│   │   ├── messages.js      # Messagerie privée
│   │   └── settings.js      # Paramètres du compte
│   ├── home.html
│   ├── explore.html
│   ├── profile.html
│   ├── messages.html
│   └── settings.html
└── Dockerfile
```

## Installation

### Prérequis
- Go 1.25+
- Serveur MySQL accessible

### Lancement local

```bash
git clone https://github.com/leo-music/The-Hive.git
cd The-Hive
```

Créer un fichier `.env` à la racine :

```env
DATABASE_URL=mysql://user:password@host:port/the-hive
APP_ADDR=:8080
```

Puis lancer :

```bash
go mod tidy
go run .
```

L'application sera accessible sur `http://localhost:8080`.

### Docker

```bash
docker build -t the-hive .
docker run -p 8080:80 --env-file .env the-hive
```

## API Routes

### Routes publiques

| Méthode | URL | Description |
|---|---|---|
| GET | `/api/ping` | Health check |
| GET | `/api/posts` | Feed public (50 derniers posts) |
| GET | `/api/post?id=X` | Détail d'un post |
| GET | `/api/comments?post_id=X` | Commentaires d'un post |
| GET | `/api/categories` | Liste des catégories |
| GET | `/api/user?username=X` | Profil public d'un utilisateur |
| GET | `/api/user/posts?username=X` | Posts d'un utilisateur |
| GET | `/api/user/likes?username=X` | Posts likés par un utilisateur |
| POST | `/api/register` | Inscription |
| POST | `/api/login` | Connexion |

### Routes protégées (session requise)

| Méthode | URL | Description |
|---|---|---|
| GET | `/api/me` | Utilisateur connecté |
| GET | `/api/feed/following` | Feed des abonnements |
| GET | `/api/conversations` | Liste des conversations |
| GET | `/api/messages?with=X` | Messages avec un utilisateur |
| POST | `/api/logout` | Déconnexion |
| POST | `/api/posts` | Créer un post |
| POST | `/api/comments` | Ajouter un commentaire |
| POST | `/api/react` | Voter (upvote/downvote) |
| POST | `/api/follow` | Suivre / ne plus suivre |
| PUT | `/api/user/update` | Modifier le profil |
| POST | `/api/user/avatar` | Changer l'avatar |

### WebSocket

| URL | Description |
|---|---|
| GET | `/api/ws` | Connexion WebSocket pour la messagerie temps réel |

## Base de données

La base MySQL est initialisée automatiquement au premier lancement. Les tables créées :

- `users` — Comptes utilisateurs
- `sessions` — Sessions d'authentification (UUID + expiration)
- `categories` — Catégories de posts (pré-remplies)
- `posts` — Publications
- `comments` — Commentaires
- `reactions` — Votes (upvote/downvote)
- `followers` — Relations d'abonnement
- `messages` — Messages privés

## Déploiement

Le projet est conçu pour être déployé via **Docker** sur **Coolify** ou tout autre hébergeur supportant les conteneurs.

1. Configurer la variable `DATABASE_URL` pointant vers votre serveur MySQL.
2. Builder et lancer le conteneur Docker.
3. Le port 80 est exposé par défaut dans le Dockerfile.
4. Optionnel : placer derrière un reverse proxy (Nginx/Caddy) avec HTTPS.
