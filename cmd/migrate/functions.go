package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	feather_commons_config "github.com/guidomantilla/go-feather-commons/pkg/config"
	feather_commons_environment "github.com/guidomantilla/go-feather-commons/pkg/environment"
	feather_commons_log "github.com/guidomantilla/go-feather-commons/pkg/log"
	"github.com/spf13/cobra"

	"github.com/guidomantilla/vaccination-record-system/pkg/config"
)

func UpCmdFn(_ *cobra.Command, args []string) {
	var err error
	err = handleMigration(args, func(migration *migrate.Migrate) error {

		if err = migration.Up(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err.Error())
	}
}

type MigrationFunction func(migration *migrate.Migrate) error

func handleMigration(args []string, fn MigrationFunction) error {

	var err error
	ctx := context.Background()

	osArgs := os.Environ()
	environment := feather_commons_environment.NewDefaultEnvironment(feather_commons_environment.WithArrays(osArgs, args))

	var cfg config.Config
	if err = feather_commons_config.Process(ctx, environment, &cfg); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	url := strings.Replace(*cfg.DatasourceUrl, ":username", *cfg.DatasourceUsername, 1)
	url = strings.Replace(url, ":password", *cfg.DatasourcePassword, 1)
	url = strings.Replace(url, ":server", *cfg.DatasourceServer, 1)
	url = strings.Replace(url, ":service", *cfg.DatasourceService, 1)

	var db *sql.DB
	if db, err = sql.Open("mysql", url); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	workingDirectory, _ := os.Getwd()
	feather_commons_log.Info(fmt.Sprintf("working directory: %s", workingDirectory))
	migrationsDirectory := filepath.Join(workingDirectory, "resources/migrations/mysql")

	var driver database.Driver
	if driver, err = mysql.WithInstance(db, &mysql.Config{}); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	var migration *migrate.Migrate
	if migration, err = migrate.NewWithDatabaseInstance("file:///"+migrationsDirectory, *cfg.DatasourceService, driver); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	if err = fn(migration); err != nil {
		feather_commons_log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	return nil
}
