package main

import (
	"fmt"
	"log/slog"
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
		"--config", 
		"./internal/store/pgstore/migrations/tern.conf",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		slog.Error("Erro ao realizar migração", "erro", string(out))
		fmt.Println(string(out))
		panic(err)
	}

}