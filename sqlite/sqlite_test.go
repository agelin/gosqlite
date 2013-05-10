package sqlite

import (
	"fmt"
	"testing"
)

func TestExecFirstAsString(t *testing.T) {
    defer func(){
        if r := recover(); r != nil {
            t.Error(fmt.Sprintf("Failed to exec first as string: %s", r))
        }
    }()
	db, _ := Open(":memory:")
	db.Throwaway("CREATE TABLE test (col)")
	db.Throwaway("INSERT INTO test VALUES ('value')")
	result, err := db.ExecFirstAsString("SELECT * FROM test")
	if err != nil || result != "value" {
		t.Error("Failed to exec first as string")
	}
}

func TestExecFirstAsInt(t *testing.T) {
    defer func(){
        if r := recover(); r != nil {
            t.Error(fmt.Sprintf("Failed to exec first as string: %s", r))
        }
    }()
	db, _ := Open(":memory:")
	db.Throwaway("CREATE TABLE test (col)")
	db.Throwaway("INSERT INTO test VALUES (1)")
	result, err := db.ExecFirstAsInt("SELECT * FROM test")
	if err != nil || result != 1 {
		t.Error("Failed to exec first as int")
	}
}

func TestFinalizeNilStmnt(t *testing.T) {
	var stmnt *Stmt = nil
	err := stmnt.Finalize()
	if err == nil || err.Error() != "Finalize called on nil Statement" {
		t.Error("Failed to error or provide a useful error message when Finalize was called on a nil Statement")
	}
}

func TestSQLNilStmnt(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "SQL called on nil Statement" {
			t.Error("Failed to panic or provide a useful error message when SQL was called on nil Statement")
		}
	}()
	var stmnt *Stmt = nil
	stmnt.SQL()
}

func TestNanosecondsNilStmnt(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "Nanoseconds called on nil Statement" {
			t.Error("Failed to panic or provide a useful error message when Nanoseconds was called on nil Statement")
		}
	}()
	var stmnt *Stmt
	stmnt.Nanoseconds()
}

func TestSafeExecToStrings(t *testing.T) {
    defer func(){
        if r := recover(); r != nil {
            t.Error(fmt.Sprintf("Failed to exec first as string: %s", r))
        }
    }()
	db, _ := Open(":memory:")
	db.Throwaway("CREATE TABLE test(col)")
	db.Throwaway("INSERT INTO test VALUES ('works')")
	_, delete_err := db.SafeExecToStrings("DELETE FROM test")
	if delete_err != nil {
		t.Error(fmt.Sprintf("Failed to safely execute delete: %s", delete_err))
	}
	result, err := db.ExecToStrings("SELECT * FROM test")
	if err != nil || result[0][0] != "works" {
		t.Error(fmt.Sprintf("Failed to execute select after safely executing delete: %s", err))
	}
}

func TestDropAllTables(t *testing.T) {
    defer func(){
        if r := recover(); r != nil {
            t.Error(fmt.Sprintf("Failed to drop all tables: %s", r))
        }
    }()
	db, _ := Open(":memory:")
	db.Throwaway("CREATE TABLE test(col)")
	pre_result, pre_result_err := db.ExecToStrings("SELECT tbl_name FROM sqlite_master")
	if pre_result_err != nil || pre_result[0][0] != "test" {
		t.Error("Failed to create test table and select it from sqlite_master")
	}
	db.DropAllTables()
	post_result, post_result_err := db.ExecToStrings("SELECT tbl_name FROM sqlite_master")
	if post_result_err != nil || len(post_result) != 0 {
		t.Error(fmt.Sprintf("Tables that should have been dropped weren't: %#v", post_result))
	}
}
