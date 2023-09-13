package database

import (
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "go-mangosteen"
	port     = 5432
	user     = "mangosteen"
	passrord = "123456"
	dbname   = "mangosteen_dev"
)

// func Connect() {
// 	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, passrord, dbname)
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	DB = db
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.Println("Connected to database")
// }

func Close() {
	DB.Close()
	log.Println("Close database connection")
}
func CreateTables() {
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

// func Migrate() {
// 	_, err := DB.Exec(
// 		`ALTER TABLE users ADD COLUMN address VARCHAR(200);`,
// 	)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.Println("Migrate table users success")
// }
