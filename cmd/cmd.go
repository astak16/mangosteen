package cmd

import (
	"fmt"
	"mangosteen/internal/database"
	"mangosteen/internal/email"
	"mangosteen/internal/router"

	"github.com/spf13/cobra"
)

func Run() {
	rootCmd := &cobra.Command{
		Use: "mangosteen",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	srvCmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			RunServer()
		},
	}
	dbCmd := &cobra.Command{
		Use: "db",
	}

	createMigrate := &cobra.Command{
		Use: "create:migration",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			database.CreateMigration(args[0])
		},
	}

	migrateCmd := &cobra.Command{
		Use: "migration:up",
		Run: func(cmd *cobra.Command, args []string) {
			database.Migrate()
		},
	}

	migrateDownCmd := &cobra.Command{
		Use: "migration:down",
		Run: func(cmd *cobra.Command, args []string) {
			database.MigrateDown()
		},
	}

	emailCmd := &cobra.Command{
		Use: "email",
		Run: func(cmd *cobra.Command, args []string) {
			email.Send()
		},
	}
	database.Connect()
	defer database.Close()

	rootCmd.AddCommand(srvCmd, dbCmd, emailCmd)
	dbCmd.AddCommand(createMigrate, migrateCmd, migrateDownCmd)

	rootCmd.Execute()

}

func RunServer() {
	r := router.New()

	r.Run(":8080")
}
