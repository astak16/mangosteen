package database

import (
	"database/sql"
	"fmt"
	"log"
	"mangosteen/sql/queries"
	"os"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect() {
	if DB != nil {
		return
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, passrord, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	DB = db
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to database")
}

func NewQuery() *queries.Queries {
	return queries.New(DB)
}

func CreateMigration(filename string) {
	fmt.Println(filename, "uccs")
	err := exec.Command("migrate", "create", "-ext", "sql", "-dir", "sql/migrations", filename).Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func Migrate() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s/sql/migrations", pwd),
		"postgres://mangosteen:123456@go-mangosteen:5432/mangosteen_dev?sslmode=disable",
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatalln(err)
	}
}

func MigrateDown() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s/sql/migrations", pwd),
		"postgres://mangosteen:123456@go-mangosteen:5432/mangosteen_dev?sslmode=disable",
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = m.Steps(-1)
	if err != nil {
		log.Fatalln(err)
	}
}
