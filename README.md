# 🐝 The Hive - Forum Project

Bienvenue sur le dépôt de **The Hive**, un projet de création d'un forum en ligne classique en partant de zéro. Ce projet est conçu pour être déployable sur un VPS et respecte des contraintes techniques strictes d'architecture.

## 🛠️ Stack Technique & Contraintes

Ce projet est développé avec les technologies suivantes pour garantir légèreté et contrôle total :
* **Backend :** Serveur web en Golang pur.
* **Base de données :** Gérée et administrée avec SQLite (`sqlite3`).
* **Frontend :** HTML5, CSS3, JS pur (Vanilla).
* **Dépendances autorisées :** Uniquement les packages standards de Go, complétés par :
    * `bcrypt` : Pour le hashage sécurisé des mots de passe.
    * `uuid` : Pour la gestion des sessions de connexion via cookies (avec temps d'expiration).
    * `go-sqlite3` : Driver pour la base de données SQLite.
* **Architecture de navigation :** Le site respecte la règle stricte d'une URL par page.

## 🗺️ Architecture du Site (Routage)

L'application est divisée entre un espace public accessible à tous et un espace membre sécurisé.

### 🌍 Espace Public (Visiteurs non connectés)
L'utilisateur non connecté peut lire les sujets, posts et commentaires.
* `GET /` : Page d'accueil (Liste des catégories et sujets récents).
* `GET /category?id=X` : Liste des sujets liés à une catégorie spécifique.
* `GET /thread?id=Y` : Lecture d'un sujet, de son post principal et de ses commentaires associés.
* `GET /login` : Page contenant le formulaire de connexion.
* `GET /register` : Page contenant le formulaire d'inscription.

### 🔒 Espace Membre (Utilisateurs connectés)
L'utilisateur connecté possède les droits de création et d'interaction.
* `GET /thread/create` : Page pour créer un nouveau sujet (lié à une catégorie).
* `GET /dashboard` : Espace personnel intégrant le système de filtrage pour afficher les sujets likés ou postés par l'utilisateur.

### ⚙️ Routes d'Action (Endpoints API internes)
Ces routes ne retournent pas de pages HTML mais traitent les formulaires (méthode `POST`) :
* `POST /auth/register` : Traite l'inscription et hashe le mot de passe.
* `POST /auth/login` : Vérifie les identifiants et génère le cookie de session `uuid`.
* `POST /auth/logout` : Détruit le cookie de session.
* `POST /comment/add` : Ajoute un commentaire à un post existant.
* `POST /react` : Gère les likes et dislikes sur les posts et commentaires.

## 🗄️ Structure de la Base de Données

La base de données relationnelle SQLite est structurée autour de 5 tables principales :
1.  **Users :** Stocke les identifiants, emails et mots de passe hashés.
2.  **Categories :** Définit les sections du forum pour le système de filtrage.
3.  **Threads :** Les sujets de discussion, liés à un utilisateur et une catégorie.
4.  **Comments :** Les réponses apportées aux sujets.
5.  **Reactions :** Table de liaison gérant les Likes/Dislikes sur les threads et les comments.

## 👥 Équipe et Répartition

Le projet est divisé en trois pôles d'expertise :
* **Dev 1 (Socle BDD) :** Conception du schéma SQLite, requêtes SQL et structures Go.
* **Dev 2 (Core Backend) :** Serveur Go, routage, authentification (`bcrypt`) et sessions (`uuid`).
* **Dev 3 (Intégration Frontend) :** Templates HTML/CSS, formulaires et rendu visuel.