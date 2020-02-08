package models

import (
	"database/sql"
	"fmt"

	log "github.com/go-pkgz/lgr"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Setting for Db
type DbSettings struct {
	Port     int
	Host     string
	User     string
	Password string
	Dbname   string
}

// InitDB
func InitDB(dbSettings DbSettings) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbSettings.Host, dbSettings.Port, dbSettings.User, dbSettings.Password, dbSettings.Dbname)

	log.Printf("[DEBUG] Starting database %s", psqlInfo)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Ping to database")
	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
