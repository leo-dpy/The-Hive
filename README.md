# 🐝 The Hive - Forum Project

Forum classique développé avec un serveur web en **Golang pur**, une base **SQLite**, des templates **HTML/CSS** et du **JavaScript Vanilla**.

## Fonctionnalités incluses

- Consultation publique des catégories, sujets et commentaires.
- Inscription et connexion utilisateur.
- Hashage des mots de passe avec `bcrypt`.
- Sessions par cookie `uuid` avec expiration de 7 jours.
- Création de sujets liés à une catégorie.
- Ajout de commentaires sur un sujet.
- Likes/dislikes sur les sujets et les commentaires.
- Dashboard membre avec filtrage : sujets postés et sujets likés.
- Une URL par page : `/`, `/category`, `/thread`, `/login`, `/register`, `/thread/create`, `/dashboard`.

## Stack et dépendances

Packages externes utilisés uniquement :

```bash
golang.org/x/crypto/bcrypt
github.com/google/uuid
modernc.org/sqlite
```

Le reste utilise la bibliothèque standard Go : `net/http`, `html/template`, `database/sql`, etc.

## Installation

```bash
git clone <ton-repo> the-hive
cd the-hive
go mod tidy
go run .
```

L'application démarre par défaut sur :

```text
http://localhost:8080
```

Tu peux changer le port avec :

```bash
APP_ADDR=:3000 go run .
```

## Base de données

La base SQLite est créée automatiquement au lancement dans :

```text
data/hive.db
```

Les catégories par défaut sont insérées automatiquement si la table `categories` est vide.

### Tables principales

- `users`
- `categories`
- `threads`
- `comments`
- `reactions`

Une table technique `sessions` est aussi utilisée pour conserver les sessions UUID et leur expiration.

## Routes

### Pages publiques

| Méthode | URL | Rôle |
|---|---|---|
| GET | `/` | Accueil, catégories, sujets récents |
| GET | `/category?id=X` | Sujets d'une catégorie |
| GET | `/thread?id=Y` | Sujet, contenu, commentaires |
| GET | `/login` | Page de connexion |
| GET | `/register` | Page d'inscription |

### Pages membres

| Méthode | URL | Rôle |
|---|---|---|
| GET | `/thread/create` | Formulaire de création de sujet |
| POST | `/thread/create` | Création du sujet |
| GET | `/dashboard` | Sujets postés ou likés |

### Actions internes

| Méthode | URL | Rôle |
|---|---|---|
| POST | `/auth/register` | Inscription |
| POST | `/auth/login` | Connexion |
| POST | `/auth/logout` | Déconnexion |
| POST | `/comment/add` | Ajout de commentaire |
| POST | `/react` | Like/dislike |

## Lancer sur un VPS

1. Installer Go. SQLite est utilisé via le driver Go, donc aucune installation SQLite séparée n’est nécessaire pour lancer le projet.
2. Copier le dossier du projet sur le serveur.
3. Exécuter `go mod tidy` puis `go build -o the-hive`.
4. Lancer le binaire : `APP_ADDR=:8080 ./the-hive`.
5. Option conseillé : placer l'application derrière Nginx avec HTTPS.

> Note : cette version utilise un driver SQLite sans CGO, donc aucun GCC/MSYS2 n’est nécessaire.

## Variante Windows sans MSYS2 / sans GCC

Cette version utilise `modernc.org/sqlite`, un driver SQLite compatible avec `database/sql` qui ne nécessite pas CGO.

Commandes sous PowerShell :

```powershell
go env -w CGO_ENABLED=0
go mod tidy
go run .
```

Puis ouvrir : http://localhost:8080

> Attention : cette variante remplace `github.com/mattn/go-sqlite3` par `modernc.org/sqlite` pour éviter l'installation de GCC/MSYS2 sur Windows. Si votre correction impose strictement `go-sqlite3`, il faudra reprendre la version CGO.
