package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB() error {
	db, err := sql.Open("sqlite3", "./calendars.db")
	stmt := "CREATE TABLE IF NOT EXISTS calendars(id INTEGER PRIMARY KEY, url TEXT, replacementSummary TEXT);"
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ReadRecord(id string) (string, error) {
	db, err := sql.Open("sqlite3", "./calendars.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	if err != nil {
		return "", err
	}

	stmt, err := db.Prepare("select url from calendars where id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var url string
	err = stmt.QueryRow("1").Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil

}

func WriteRecord(params CalendarParams) {
	db, err := sql.Open("sqlite3", "./calendars.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt := "INSERT INTO calendars(url, replacementSummary) VALUES(?, ?);"
	result, err := db.Exec(stmt, params.Url, params.ReplacementSummary)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result.LastInsertId())
}
