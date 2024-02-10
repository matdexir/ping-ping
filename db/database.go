package db

import (
	// "context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type PostDB struct {
	Database *sql.DB
}

func (p *PostDB) CreateTable() error {
	statement, err := p.Database.Prepare(`
    CREATE TABLE IF NOT EXISTS people(
      id INTEGER PRIMARY KEY, 
      title TEXT, startAt TEXT, 
      endAt TEXT, 
      ageStart INTEGER, 
      ageEnd INTEGER, 
      targetGender TEXT, 
      targetCountries TEXT, 
      targetPlatforms TEXT )`)

	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to prepare statement")
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to create table")
		return err
	}

	return nil
}
