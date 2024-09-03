FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./


RUN go install github.com/air-verse/air@latest

COPY . .
EXPOSE 3000

CMD ["air"]