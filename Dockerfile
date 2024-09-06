FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./


RUN go install github.com/air-verse/air@latest



COPY . .

RUN go install github.com/jackc/tern/v2@latest  

RUN tern migrate -m ./internal/store/pgstore/migrations --config ./internal/store/pgstore/migrations/tern.conf

EXPOSE 3000

CMD ["air"]