package db

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Record struct {
	Url                string
	ReplacementSummary string
}

// TODO: write a class that stores dbFilename
func InitDB(dbFilename string) error {
	db, err := sql.Open("sqlite3", dbFilename)
	stmt := "CREATE TABLE IF NOT EXISTS calendars(" +
		"id TEXT PRIMARY KEY, " +
		"url TEXT, " +
		"replacementSummary TEXT);"
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ReadRecord(dbFilename string, id string) (Record, error) {
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return Record{}, err
	}
	defer db.Close()

	if err != nil {
		return Record{}, err
	}

	stmt, err := db.Prepare("select url, replacementSummary from calendars where id = ?")
	if err != nil {
		return Record{}, err
	}
	defer stmt.Close()

	var url, replacementSummary string
	err = stmt.QueryRow(id).Scan(&url, &replacementSummary)
	if err == sql.ErrNoRows {
		return Record{}, err
	} else if err != nil {
		return Record{}, err
	}

	return Record{url, replacementSummary}, nil

}

func WriteRecord(dbFilename string, record Record) (string, error) {
	id := uuid.New().String()
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt := "INSERT INTO calendars(id, url, replacementSummary) VALUES(?, ?, ?);"
	_, err = db.Exec(stmt, id, record.Url, record.ReplacementSummary)
	if err != nil {
		return "", err
	}

	return id, nil

}
