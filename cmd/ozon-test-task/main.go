package main

import (
	"database/sql"
	"os"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/adaptor"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/migrations"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/service/http"

	_ "github.com/lib/pq"
)

func main() {
	address := os.Getenv("HTTP_ADDRESS")
	if address == "" {
		address = ":80"
	}

	postgresCreds := os.Getenv("POSTGRES")

	var repo adaptor.Repository
	if postgresCreds == "" {
		repo = adaptor.NewInMemoryRepository()
	} else {
		db, err := sql.Open("postgres", postgresCreds)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = migrations.RunMigrations(db)
		if err != nil {
			panic(err)
		}

		repo = adaptor.NewPostgresRepository(db)
	}

	err := http.Run(address, repo)
	if err != nil {
		panic(err)
	}
}
