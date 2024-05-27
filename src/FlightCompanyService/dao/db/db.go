package dbClient

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var dbInstance *sql.DB

func createClient() (*sql.DB, error) {

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("FLIGHTCOMPANY_DB_USER"),
		Passwd: os.Getenv("FLIGHTCOMPANY_DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("FLIGHTCOMPANY_DB_HOST"),
		DBName: os.Getenv("FLIGHTCOMPANY_DB_NAME"),
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())

	return db, err
}

func InitDB() {
	var err error
	dbInstance, err = createClient()
	if err == nil {
		pingErr := dbInstance.Ping()
		if pingErr != nil {
			err = pingErr
			fmt.Printf("DB Connection failed: %+v\n", pingErr)
			panic(err)
		} else {
			fmt.Println("DB Connection Succedeed")
		}
	}
}

func CloseClient() {
	log.Println("closing DB client")
	_ = dbInstance.Close()
}

func GetInstance() *sql.DB {
	return dbInstance
}
