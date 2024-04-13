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
		User:   os.Getenv("ACMESKY_DB_USER"),
		Passwd: os.Getenv("ACMESKY_DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("ACMESKY_DB_HOST"),
		DBName: os.Getenv("ACMESKY_DB_NAME"),
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
		} else {
			fmt.Println("DB Connection Succedeed")
		}
	}
}

func closeClient(client *sql.DB) {
	log.Println("closing DB client")
	_ = client.Close()
}

func GetInstance() *sql.DB {
	return dbInstance
}
