package cmd

import (
	"fmt"
	"log"
	"mangosteen/global"
	"mangosteen/internal/database"
	"mangosteen/internal/email"
	"mangosteen/internal/jwt_helper"
	"mangosteen/internal/router"
	"net/http"
	"os"
	"os/exec"

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

	generateHMACKeyCmd := &cobra.Command{
		Use: "generateHMACKey",
		Run: func(cmd *cobra.Command, args []string) {
			bytes, _ := jwt_helper.GenerateHMACKey()
			fmt.Println(global.JwtPath)
			if err := os.WriteFile(global.JwtPath, bytes, 0644); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("HMAC key generated")
		},
	}

	coverCmd := &cobra.Command{
		Use: "cover",
		Run: func(cmd *cobra.Command, args []string) {
			os.MkdirAll("coverage", os.ModePerm)
			if err := exec.Command("MailHog").Start(); err != nil {
				log.Println(err)
			}
			if err := exec.Command("go", "test", "-coverprofile=coverage/coverage.out", "./...").Run(); err != nil {
				log.Fatalln(err)
			}
			if err := exec.Command("go", "tool", "cover", "-html=coverage/coverage.out", "-o", "coverage/index.html").Run(); err != nil {
				log.Fatalln(err)
			}
			port := "8888"
			if len(args) > 0 {
				port = args[0]
			}
			fmt.Printf("start server: http://localhost:%s/coverage/index.html\n", port)
			if err := http.ListenAndServe(":"+port, http.FileServer(http.Dir("."))); err != nil {
				log.Fatalln(err)
			}
		},
	}

	database.Connect()
	defer database.Close()

	rootCmd.AddCommand(srvCmd, dbCmd, emailCmd, generateHMACKeyCmd, coverCmd)
	dbCmd.AddCommand(createMigrate, migrateCmd, migrateDownCmd)

	rootCmd.Execute()

}

func RunServer() {
	r := router.New()

	r.Run(":8080")
}
