package db

import "testing"

func TestReadWrite(t *testing.T) {

	db := "test.db"
	err := InitDB(db)
	if err != nil {
		t.Error(err)
	}

	expected := Record{
		Url:                "abc",
		ReplacementSummary: "123",
	}

	record := Record{
		Url:                "abc",
		ReplacementSummary: "123",
	}
	id, err := WriteRecord(db, record)
	if err != nil {
		t.Error(err)
	}
	got, err := ReadRecord(db, id)
	if err != nil {
		t.Error(err)
	}

	if got != expected {
		t.Errorf("Error:\nExpected: %s\nGot: %s", expected, got)
	}

}
