package db

import (
	// "context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type PostDB struct {
	Database *sql.DB
}

func CreateConnection() (*PostDB, error) {
	// connStr := "postgresql://postgres:postgres@127.0.0.1:8001/posts?sslmode=disable"
	// connStr := "host=0.0.0.0 port=8001 user=postgres password=postgres dbname=posts sslmode=disable"
	sql, err := sql.Open("sqlite3", "./db/file.db")
	// sql, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Unable to open database")
		return nil, err
	}

	if err = sql.Ping(); err != nil {
		log.Printf("Unable to ping the database: %v\n", err)
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
