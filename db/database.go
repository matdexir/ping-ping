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

func CreateConnection() (*PostDB, error) {
	sql, err := sql.Open("sqlite3", "./db/file.db")
	if err != nil {
		fmt.Println("Unable to open database")
		return nil, err
	}
	return &PostDB{Database: sql}, nil
}

func (pdb *PostDB) Close() error {
	return pdb.Database.Close()
}

func (pdb *PostDB) CreateTable() error {

	statement, err := pdb.Database.Prepare(`
    CREATE TABLE IF NOT EXISTS posts(
      id INTEGER PRIMARY KEY NOT NULL, 
      title TEXT, 
      startAt TEXT NOT NULL,
      endAt TEXT NOT NULL, 
      ageStart INTEGER, 
      ageEnd INTEGER, 
      targetGender TEXT, 
      targetCountry TEXT, 
      targetPlatform TEXT )`)

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
