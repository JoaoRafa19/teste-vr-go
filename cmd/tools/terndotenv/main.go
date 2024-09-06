package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Arquivo .env não encontrado, utilizando variáveis de ambiente do sistema")
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"-m",
		"./internal/store/pgstore/migrations",
		fmt.Sprintf("--host %s", os.Getenv("DATABASE_HOST")),
		fmt.Sprintf("--database %s", os.Getenv("DATABASE_NAME")),
		fmt.Sprintf("--password %s", os.Getenv("DATABASE_PASSWORD")),
		fmt.Sprintf("--port %s", os.Getenv("DATABASE_PORT")), 
		fmt.Sprintf("--user %s", os.Getenv("DATABASE_USER")),

	)

	if out, err := cmd.CombinedOutput(); err != nil {
		slog.Error("Erro ao realizar migração", "erro", string(out))
		fmt.Println(string(out))
		panic(err)
	}

}