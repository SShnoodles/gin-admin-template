FROM node:22-bookworm-slim AS web-builder
WORKDIR /src/web

COPY web/package.json web/package-lock.json ./
RUN npm ci --legacy-peer-deps --no-audit --no-fund

COPY web/ ./
RUN npm run build

FROM golang:1.21-alpine AS api-builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=web-builder /src/internal/web/dist ./internal/web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gin-admin-template ./main.go

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=api-builder /out/gin-admin-template ./gin-admin-template
COPY --from=api-builder /src/config.yml ./config.yml
COPY --from=api-builder /src/docs ./docs
COPY --from=api-builder /src/locales ./locales

RUN mkdir -p logs

EXPOSE 8080

ENTRYPOINT ["./gin-admin-template"]
