# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./frontend/
WORKDIR /app/frontend
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.22-alpine AS backend
WORKDIR /app/backend
RUN apk add --no-cache git
COPY backend/go.mod backend/go.sum* ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend /app/build/static ./ui/static
RUN CGO_ENABLED=0 go build -o /build/bin/ttrpg-toolkit ./cmd

# Stage 3: Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend /build/bin/ttrpg-toolkit /app/ttrpg-toolkit
EXPOSE 8080
CMD ["/app/ttrpg-toolkit"]
