package dbClient

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/go-sql-driver/mysql"
)

var lockDBInstance = &sync.Mutex{}
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

func closeClient(client *sql.DB) {
	log.Println("closing DB client")
	_ = client.Close()
}

func GetInstance() (*sql.DB, error) {
	var err error

	if dbInstance == nil {
		lockDBInstance.Lock()
		defer lockDBInstance.Unlock()
		if dbInstance == nil {
			dbInstance, err = createClient()

			if err != nil {
				pingErr := dbInstance.Ping()
				if pingErr != nil {
					err = pingErr
				}
			}
		}
	}

	return dbInstance, err
}
