FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./


RUN go install github.com/air-verse/air@latest



COPY . .

RUN go install github.com/jackc/tern/v2@latest  

RUN go run ./cmd/tools/terndotenv/main.go

EXPOSE 3000

CMD ["air"]