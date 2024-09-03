FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

EXPOSE 3000

CMD go run 