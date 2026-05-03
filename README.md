# 🐝 The Hive - Réseau Social & Forum

Bienvenue sur le dépôt de **The Hive**. Ce projet a pour ambition de créer une plateforme hybride, à la croisée des chemins entre un forum classique et un réseau social moderne. Construit de zéro, ce projet est conçu pour être déployable sur un VPS tout en respectant des contraintes techniques strictes d'architecture.

## 🚀 Vision du Projet

The Hive n'est pas qu'un simple forum, c'est un véritable **pseudo réseau social** centré sur les intéractions et la découverte, tout en respectant scrupuleusement les fondations d'un forum classique. Les fonctionnalités principales incluent :

*   **Catégories & Sujets :** L'organisation centrale. Les utilisateurs parcourent une multitude de **Catégories** et peuvent y créer ou lire des **Sujets (Posts)** et leurs commentaires.
*   **Système de Filtrage Avancé & Fil d'Actualité :** Un flux dynamique et un système de filtre permettant de trier les sujets **par catégorie**, mais aussi de retrouver facilement **les sujets que l'utilisateur a likés ou postés**.
*   **Interactions (Likes & Dislikes) :** Un système complet de réactions permettant de liker ou disliker les posts et commentaires pour animer la communauté.
*   **Système Social (Amis & Abonnements) :** Possibilité de suivre d'autres utilisateurs (système de "follow" ou d'amis) pour personnaliser son fil d'actualité.
*   **Messagerie Privée (MP) :** Un espace dédié pour discuter en privé entre membres (1 to 1).

## 🛠️ Stack Technique & Contraintes

Ce projet est développé avec les technologies suivantes pour garantir légèreté et contrôle total :
*   **Backend :** Serveur web en Golang pur.
*   **Base de données :** Gérée et administrée avec SQLite (`sqlite3`).
*   **Frontend :** HTML5, CSS3, JS pur (Vanilla).
*   **Dépendances autorisées :** Uniquement les packages standards de Go, complétés par :
    *   `bcrypt` : Pour le hashage sécurisé des mots de passe.
    *   `uuid` : Pour la gestion des sessions de connexion via cookies (avec temps d'expiration).
    *   `go-sqlite3` : Driver pour la base de données SQLite.
*   **Architecture de navigation :** Le site respecte la règle stricte d'une URL par page (ou d'une architecture claire pour le rendu des vues).

## 🗺️ Architecture du Site (Routage Prévisionnel)

L'application est divisée entre un espace public et un espace membre offrant des fonctionnalités sociales avancées.

### 🌍 Espace Public (Visiteurs non connectés)
*   `GET /` : Page d'accueil / Fil d'actualité global (Tendances, Catégories populaires).
*   `GET /category?id=X` : Vue d'une catégorie et liste de ses sujets.
*   `GET /subject?id=Y` : Lecture d'un sujet (post) spécifique et de ses commentaires.
*   `GET /login` & `GET /register` : Pages d'authentification.

### 🔒 Espace Membre (Utilisateurs connectés)
*   `GET /feed` : Fil d'actualité personnalisé (Sujets des amis, sujets likés/postés par l'utilisateur).
*   `GET /subject/create` : Page pour créer un nouveau sujet dans une catégorie.
*   `GET /profile?user=Z` : Profil d'un utilisateur (permettant de s'abonner / l'ajouter en ami).
*   `GET /messages` : Boîte de réception de la messagerie privée (MP).
*   `GET /messages?user=W` : Fil de discussion privé avec un utilisateur spécifique.

### ⚙️ Routes d'Action (Endpoints API internes)
*   `POST /auth/register`, `POST /auth/login`, `POST /auth/logout` : Gestion de l'authentification.
*   `POST /follow` : S'abonner / Ajouter un ami.
*   `POST /react` : Liker ou disliker un sujet ou un commentaire.
*   `POST /comment/add` : Ajouter un commentaire à un sujet.
*   `POST /message/send` : Envoyer un message privé (MP).

## 🗄️ Structure de la Base de Données (Évolution)

La base de données relationnelle SQLite sera structurée pour soutenir cet aspect réseau social tout en respectant le format forum :
1.  **Users :** Identifiants, emails, mots de passe hashés, infos de profil.
2.  **Follows/Friends :** Table de liaison pour gérer le système d'abonnements ou d'amis entre utilisateurs.
3.  **Categories :** Les différentes sections thématiques du forum.
4.  **Subjects :** Les sujets (posts principaux) créés par les utilisateurs, rattachés à une *Category*.
5.  **Comments :** Les réponses publiées au sein d'un *Subject*.
6.  **Reactions :** Table de liaison gérant les Likes et Dislikes sur les *Subjects* et *Comments*.
7.  **PrivateMessages :** Table stockant les messages privés (MP) envoyés entre deux utilisateurs (`sender_id`, `receiver_id`).

## 🌟 Fonctionnalités Bonus (100% Go Natif)

Pour pousser l'aspect "Réseau Social" plus loin tout en respectant strictement l'absence de framework (pur Go et Vanilla JS), le projet intègre ces fonctionnalités avancées :

*   **🔔 Notifications en Temps Réel (SSE) :** Un système d'alertes en direct (likes, mentions, MP, nouveaux followers) sans utiliser WebSockets, en exploitant nativement les *Server-Sent Events* via l'interface `http.Flusher` de la librairie standard Go.
*   **🏷️ Mentions et Tags (@utilisateur) :** Le parsing des messages côté serveur (via `regexp`) permet de détecter les `@pseudo`, de les transformer en liens cliquables vers le profil cible, et de déclencher une notification.
*   **🔍 Moteur de Recherche Global (FTS) :** Une recherche ultra-rapide sur tout le site (utilisateurs, sujets, discussions) exploitant le module **Full-Text Search** natif de SQLite (`go-sqlite3`).
*   **🖼️ Personnalisation du Profil (Avatars) :** Gestion de l'upload d'images de profil grâce à `r.ParseMultipartForm` de Go, avec vérification sécurisée des fichiers et stockage local.

## 👥 Équipe et Répartition

Le projet est divisé en trois pôles d'expertise (à titre indicatif) :
*   **Dev 1 (Socle BDD & Logique Sociale) :** Schéma SQLite complexe (relations amis, requêtes de flux d'actualité), structures Go.
*   **Dev 2 (Core Backend) :** Serveur Go, routage, authentification (`bcrypt`), sessions (`uuid`) et gestion des MP.
*   **Dev 3 (Intégration Frontend & UI) :** Templates HTML/CSS, UI d'un réseau social (fil d'actu fluide, chat box), interactions JavaScript Vanilla.
