# Go image
FROM golang:1.24-alpine

# Ish papkasi
WORKDIR /app

# Go modules nusxalash
COPY go.mod go.sum ./
RUN go mod download

# Kodni nusxalash
COPY internal/ ./internal
COPY cmd/ ./cmd

# config.yaml ni WORKDIR ichiga nusxalash
COPY config.yaml ./config.yaml

# Build
RUN go build -o main ./cmd/main.go

# Portni ochish
EXPOSE 8080

# Run
CMD ["./main"]
