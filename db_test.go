package main

import "testing"

func testReadWrite(t *testing.T) {

	db := "test.db"
	err := InitDB(db)
	if err != nil {
		t.Error(err)
	}
	expected := "hello there\nnice to meet you"
	id, err := WriteRecord(db, "abc", "123", "expected")
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
