package fsb

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"io/fs"
)

const schemaVersion = "fsb_schema_version"

//go:embed migrations/*.sql
var migrationFiles embed.FS

func runDBMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	ternMigrator, err := migrate.NewMigrator(ctx, conn.Conn(), schemaVersion)
	if err != nil {
		return fmt.Errorf("failed to create tern migrator: %w", err)
	}

	migrationRoot, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	if err = ternMigrator.LoadMigrations(migrationRoot); err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if err = ternMigrator.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}
