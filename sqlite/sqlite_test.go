package sqlite

import (
	"testing"
)

func open(t *testing.T) *Conn {
  db, err := Open(":memory:")
  if err != nil {
    t.Error("Failed to open memory database")
  }
  return db
}

func (db *Conn) create(t *testing.T, sql string) {
  err := db.Exec(sql)
  if err != nil {
    t.Error("Failed to create")
  }
}

func (db *Conn) insert(t *testing.T, sql string) {
  err := db.Exec(sql)
  if err != nil {
    t.Error("Failed to insert")
  }
}

func TestExecFirstAsString(t *testing.T) {
	db := open(t)
  db.create(t, "CREATE TABLE test (col)")
  db.insert(t, "INSERT INTO test VALUES ('value')")
  result, err := db.ExecFirstAsString("SELECT * FROM test")
  if err != nil || result != "value" {
    t.Error("Failed to exec first as string")
  }
}

func TestExecFirstAsInt(t *testing.T) {
	db := open(t)
  db.create(t, "CREATE TABLE test (col)")
  db.insert(t, "INSERT INTO test VALUES (1)")
  result, err := db.ExecFirstAsInt("SELECT * FROM test")
  if err != nil || result != 1 {
    t.Error("Failed to exec first as int")
  }
}

