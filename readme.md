![header](https://capsule-render.vercel.app/api?type=venom&color=auto&height=400&section=header&text=Teste%20VR&fontSize=90&rotate=10)

![go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)![pgsql](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)

# Teste VR

# Execução 

```shell 
make up && make logs
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


## __Observação__
É necessário criar o arquivo .env na raiz do projeto e na pasta service core com as variáveis de ambiente necessárias
ex.:
```.env
DATABASE_PORT=5432
DATABASE_USER="postgres"
DATABASE_PASSWORD="123456789"
DATABASE_NAME="plataforma-ead"
DATABASE_HOST="service-core-db"
```

### Comandos principais
***Rodando o container***
```shell
make up
```
___(ou com logs)___
```shell
MODE=l make up
```
```shell
make up && make logs
```

***Logs***
>Com o container já de pé ele vai acoplar o terminal ao terminal de logs do docker.

```shell
make logs
```

Restart
> Reinicia o container
```shell
make restart
```

Parar
>Encerra a execução da aplicação
```shell
make down
```


## Dev dependencies

AIR (live reload do go)
```
go install github.com/air-verse/air@latest 
```
### Comandos do projeto
Compilar os arquivo do sqlc
```shell
go gen ./...
```


## Testes

### Service Core
Executar testes service core
```shell
make test-service-core
```