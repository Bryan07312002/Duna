package main

import (
	"fmt"
	"os"

	"duna/internal/database"

	"github.com/joho/godotenv"
)

const HELP_MESSAGE = `
Usage: duna <command> [options]

Commands:
	migrate 	Run database migrations

	rollback // TODO

	serve      Start the application server
     options:
       -port int   port to listen on (default 8080)
       -env string environment (default 'development') // TODO
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(HELP_MESSAGE)
		os.Exit(1)
	}

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	switch os.Args[1] {
	case "migrate":
		handleMigrate()
	default:
		fmt.Println("unknown command")
	}
}

func handleMigrate() {
	db, err := database.NewDatabase()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := db.Migrate(); err != nil {
		fmt.Println(err.Error())
	}
}
