# ════════════════════════════════════════════
# The Hive — Multi-stage Dockerfile
# ════════════════════════════════════════════

# ── Stage 1: Build ──
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Installer les dépendances système
RUN apk add --no-cache git

# Copier les fichiers de dépendances d'abord (cache Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le reste du code source
COPY . .

# Compiler l'application de façon statique
RUN CGO_ENABLED=0 GOOS=linux go build -o the-hive .

# ── Stage 2: Runtime ──
FROM alpine:latest

WORKDIR /app

# Certificats SSL (pour les connexions DB externes) et fuseau horaire
RUN apk add --no-cache ca-certificates tzdata

# Copier le binaire et le dossier public depuis le builder
COPY --from=builder /app/the-hive .
COPY --from=builder /app/public ./public

# Créer le dossier uploads pour les avatars (volume persistant)
RUN mkdir -p /app/public/uploads/avatars

# Définir le port d'écoute et l'exposer
EXPOSE 80

# Lancer l'application
CMD ["./the-hive"]
