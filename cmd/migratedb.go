package cmd

import (
	db_migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/wahed-tech/banking-ops-svc/conf"
	"github.com/wahed-tech/go-lib/v2/cmds"
	"github.com/wahed-tech/go-lib/v2/errors"
)

var migratedbCmd = &cobra.Command{
	Use:   "migratedb",
	Short: "migratedb migrate db to last versions",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	RunE:  cmds.WithErrorLog(migrateUp),
	Short: "apply all up database migrations",
}

func init() {
	rootCmd.AddCommand(migratedbCmd)
	migratedbCmd.AddCommand(migrateUpCmd)
}

func migrateUp(_ *cobra.Command, _ []string) error {
	m, err := db_migrate.New("file://"+conf.R.Postgres.MigrationsDir, conf.R.Postgres.URL())
	if err != nil {
		return errors.Wrap(err)
	}
	defer m.Close()

	err = m.Up()
	if errors.Is(err, db_migrate.ErrNoChange) {
		return nil
	}

	return errors.Wrap(err)
}
