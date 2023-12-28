package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/guidomantilla/vaccination-record-system/cmd/serve"
)

func ExecuteAppCmd() {

	appCmd := &cobra.Command{}
	appCmd.AddCommand(createServeCmd(), createMigrateCmd())

	if err := appCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}

func createServeCmd() *cobra.Command {

	serveCmd := &cobra.Command{
		Use: "serve",
		Run: serve.ExecuteCmdFn,
	}

	return serveCmd
}

func createMigrateCmd() *cobra.Command {

	migrateUpCmd := &cobra.Command{
		Use: "up",
		//Run: migrate.UpCmdFn,
	}

	migrateDownCmd := &cobra.Command{
		Use: "down",
		//Run: migrate.DownCmdFn,
	}

	migrateCmd := &cobra.Command{
		Use: "migrate",
	}

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)

	return migrateCmd
}
