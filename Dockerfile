# Stage 1: Build frontend
FROM node:alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/dist ./web/dist
ARG VERSION=dev
RUN go build -ldflags="-X 'github.com/mano/mucal/internal/version.Version=${VERSION}'" -o /mucal ./cmd/mucal

# Stage 3: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /mucal ./
EXPOSE 8080
ENTRYPOINT ["/app/mucal"]
CMD ["-config", "/config/config.yaml"]
