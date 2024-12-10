package postgres

import (
	"fmt"
	"music-library/config"

	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	const op = "db.postgres.Init"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// migrator := MustGetNewMigrator(MigrationsFS, migrationsDir)

	// err = migrator.ApplyMigrations(db)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

// const migrationsDir = "migrations"

// //go:embed migrations/*.sql
// var MigrationsFS embed.FS

// type Migrator struct {
// 	srcDriver source.Driver
// }

// func MustGetNewMigrator(sqlFiles embed.FS, dirName string) *Migrator {
// 	d, err := iofs.New(sqlFiles, dirName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &Migrator{
// 		srcDriver: d,
// 	}
// }

// func (m *Migrator) ApplyMigrations(db *sql.DB) error {
// 	driver, err := postgres.WithInstance(db, &postgres.Config{})
// 	if err != nil {
// 		return fmt.Errorf("unable to create db instance: %v", err)
// 	}

// 	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
// 	if err != nil {
// 		return fmt.Errorf("unable to create migration: %v", err)
// 	}

// 	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
// 		return fmt.Errorf("unable to apply migrations %v", err)
// 	}

// 	return nil
// }
