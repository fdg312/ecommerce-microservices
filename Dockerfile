FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth-server ./cmd/auth/main.go
EXPOSE 50051
CMD ["./auth-server"]
