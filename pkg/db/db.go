package db

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var calendarCache = make(map[string]string)

func ClearCache() {
	calendarCache = make(map[string]string)
}

type Record struct {
	Url                string
	ReplacementSummary string
	CalendarBody       string
}

// TODO: write a class that stores dbFilename
func InitDB(dbFilename string) error {
	db, err := sql.Open("sqlite3", dbFilename)
	stmt := "CREATE TABLE IF NOT EXISTS calendars(" +
		"id TEXT PRIMARY KEY, " +
		"url TEXT, " +
		"replacementSummary TEXT, " +
		"calendarBody TEXT" +
		");"
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ReadRecord(dbFilename string, id string) (string, error) {
	calendarBody, exists := calendarCache[id]
	if exists {
		log.Printf("Cache hit for %s\n", id)
		return calendarBody, nil
	} else {
		log.Printf("Cache miss for %s\n", id)
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return "", err
	}
	defer db.Close()

	if err != nil {
		return "", err
	}

	stmt, err := db.Prepare("select calendarBody from calendars where id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&calendarBody)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}
	calendarCache[id] = calendarBody
	return calendarBody, nil

}

func WriteRecord(dbFilename string, record Record) (string, error) {
	id := uuid.New().String()
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt := "INSERT INTO calendars(id, url, replacementSummary, calendarBody) VALUES(?, ?, ?, ?);"
	_, err = db.Exec(stmt, id, record.Url, record.ReplacementSummary, record.CalendarBody)
	if err != nil {
		return "", err
	}

	return id, nil

}
