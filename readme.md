![header](https://capsule-render.vercel.app/api?type=venom&color=auto&height=400&section=header&text=Teste%20VR&fontSize=90&rotate=10)

![go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)![pgsql](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)

# Teste VR

# Execução 

```shell 
docker compose up
```

# Desenvolvimento

Aplicação backend de uma plataforma de cursos.

Utiliza `sqlc` para gerar as interfaces das entidades das tabelas dos bancos de dados (não é um ORM) e as queries SQL.
Utiliza o `tern` para criar e executar as migations.


## Go generate

Executa os comandos declarados em `gen.go`
```go
package gen 


//go:generate go run ./cmd/tools/terndotenv/main.go
//go:generate sqlc generate -f ./internal/store/pgstore/sqlc.yml
```
```shell
go generate ./...
```

## Migrations
Utiliazando o tern para criar migrações, mas para executar com o ambiente local do docker pelo arquivo .env
utiliza o `os\exec` do go para rodar comandos no ambiente

```shell
go run ./cmd/tools/terndotenv/main.go
```

## Queries

Usa `sqlc` para gerar as queries

```shell
sqlc generate -f ./internal/store/pgstore/sqlc.yml
```


## Deps


#### Install all deps:
```shell
go mod tidy
```

- **tern**
```shell
 go install github.com/jackc/tern/v2@latest
 ```

- **sqlc**
```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
