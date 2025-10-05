FROM golang:1.25.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/githubnemo/CompileDaemon@latest
EXPOSE 8080
ENTRYPOINT CompileDaemon --build="go build -o build/goapp" -command="./build/goapp" -build-dir=/app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .