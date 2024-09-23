package adaptor_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/adaptor"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/migrations"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository(t *testing.T) {
	t.Parallel()

	postgresCreds := os.Getenv("POSTGRES")
	if postgresCreds == "" {
		t.Skip("environment variable POSTGRES not specified")
	}

	db, err := sql.Open("postgres", postgresCreds)
	require.NoError(t, err)
	defer db.Close()

	err = migrations.RunMigrations(db)
	require.NoError(t, err)

	runSharedRepositoryTests(t, adaptor.NewPostgresRepository(db))
}
