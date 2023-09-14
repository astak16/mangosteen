package main

import (
	"mangosteen/cmd"
	"mangosteen/initialize"
)

func main() {

	initialize.InitViper()
	cmd.Run()

	// database.Connect()
	// database.Migrate()
	// defer database.Close()
	// cmd.RunServer()
}
