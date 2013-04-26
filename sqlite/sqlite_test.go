package sqlite

import (
	"testing"
)

func TestOpenMemory(t *testing.T) {
	_, open_err := Open(":memory:")
  if open_err != nil {
    t.Error("Failed to open in memory only database")
  }
}

func TestCreateTableInMemory(t *testing.T) {
	db, open_err := Open(":memory:")
  if open_err != nil {
    t.Error("Failed to open in memory only database")
  }
  create_error := db.Exec("CREATE TABLE test (col)")
  if create_error != nil {
    t.Error("Failed to create in-memory only table")
  }
}

func TestInsertInMemory(t *testing.T) {
	db, open_err := Open(":memory:")
  if open_err != nil {
    t.Error("Failed to open in memory only database")
  }
  create_error := db.Exec("CREATE TABLE test (col)")
  if create_error != nil {
    t.Error("Failed to create in-memory only table")
  }
  insert_error := db.Exec("INSERT INTO test VALUES ('value')")
  if insert_error != nil {
    t.Error("Failed to insert into memory only database")
  }
}

func TestExecFirstAsString(t *testing.T) {
	db, open_err := Open(":memory:")
  if open_err != nil {
    t.Error("Failed to open in memory only database")
  }
  create_error := db.Exec("CREATE TABLE test (col)")
  if create_error != nil {
    t.Error("Failed to create in-memory only table")
  }
  insert_error := db.Exec("INSERT INTO test VALUES ('value')")
  if insert_error != nil {
    t.Error("Failed to insert into memory only database")
  }
  result, err := db.ExecFirstAsString("SELECT * FROM test")
  if err != nil {
    t.Error("Failed to select from test")
  }
  if result != "value" {
    t.Error("Failed to exec first as string")
  }
}
