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

# Copier le reste du code
COPY . .

# Compiler le binaire
RUN CGO_ENABLED=0 GOOS=linux go build -o the-hive .

# ── Stage 2: Runtime ──
FROM alpine:3.20

WORKDIR /app

# Certificats SSL (pour les connexions DB externes)
RUN apk add --no-cache ca-certificates tzdata

# Copier le binaire compilé
COPY --from=builder /app/the-hive .

# Copier les fichiers statiques
COPY --from=builder /app/public ./public

# Créer le dossier uploads pour les avatars (volume persistant)
RUN mkdir -p /app/public/uploads/avatars

# Port exposé
EXPOSE 8080

# Lancement
CMD ["./the-hive"]
