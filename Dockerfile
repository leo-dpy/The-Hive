# ── Stage 1: Build ──
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Installer les dépendances système
RUN apk add --no-cache git

# Copier les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code source
COPY . .

# Compiler l'application
RUN CGO_ENABLED=0 GOOS=linux go build -o the-hive .

# ── Stage 2: Runtime ──
FROM alpine:latest

WORKDIR /app

# Certificats SSL et zone horaire
RUN apk add --no-cache ca-certificates tzdata

# Copier le binaire et les fichiers statiques
COPY --from=builder /app/the-hive .
COPY --from=builder /app/public ./public

# Dossier pour les uploads
RUN mkdir -p /app/public/uploads/avatars

# Configuration du port
# Note: On force l'application à écouter sur le port 80 pour Coolify
ENV APP_ADDR=":80"
EXPOSE 80

# Lancement
CMD ["./the-hive"]
