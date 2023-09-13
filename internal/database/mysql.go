package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func MySQLDatabase() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", passrord, "go-uccs-1", 3306, dbname)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	DB = db
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
}

func MySqlCreateTables() {
	sqlStatement := `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT
	);
	`
	_, err := DB.Exec(sqlStatement)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Create table users success")
}
