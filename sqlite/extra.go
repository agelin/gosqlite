// Copyright 2013 Richard B. Lyman. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sqlite

/*
#include <sqlite3.h>
*/
import "C"
import "errors"
import "fmt"

func Columns( stmnt *Stmt ) int {
    return int(C.sqlite3_column_count(stmnt.stmt))
}

// ScanAllAsString will return all rows if error is nil.
// If error is not nil, the rows successfully scanned up to
// that point are the only rows returned.
func ScanAllAsString( stmnt *Stmt ) ([][]string, error) {
    result := [][]string{}
    for stmnt.Next() {
        row, err := ScanAsString(stmnt)
        if err != nil {
            return result, err
        }
        result = append(result, row)
    }
    return result, nil
}

func ScanAsString( stmnt *Stmt ) ([]string, error) {
    result := make([]string,Columns(stmnt))
    addrs := make([]interface{},Columns(stmnt))
    for i := range result {
        addrs[i] = &result[i]
    }
    err := stmnt.Scan(addrs...)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func ScanAsInt( stmnt *Stmt ) ([]int, error) {
    result := make([]int,Columns(stmnt))
    addrs := make([]interface{},Columns(stmnt))
    for i := range result {
        addrs[i] = &result[i]
    }
    err := stmnt.Scan(addrs...)
    if err != nil {
        return nil, err
    }
    return result, nil
}

// See ScanAllAsString for an explanation of how possible errors
// affect the rows returned
func (c *Conn) ExecToStrings( sql string ) ([][]string, error) {
    stmnt, prep_err := c.Prepare(sql)
    if prep_err != nil {
        return nil, prep_err
    }
    has_rows := stmnt.Next()
    if !has_rows {
        return nil, errors.New("There are no rows in the result set")
    }
    return ScanAllAsString(stmnt)
}

func (c *Conn) SafeExecToStrings( sql string ) ([][]string, error) {
  begin_err := c.Exec("BEGIN;")
  if begin_err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to begin safe execution: %s", begin_err))
  }
  result, result_err := c.ExecToStrings(sql)
  if result_err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to execute query: %s", result_err))
  }
  rollback_err := c.Exec("ROLLBACK;")
  if rollback_err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to rollback safe execution: %s", rollback_err))
  }
  return result, nil
}

func (c *Conn) ExecFirstAsString( sql string ) (string, error) {
    stmnt, prep_err := c.Prepare(sql)
    if prep_err != nil {
        return "", prep_err
    }
    has_rows := stmnt.Next()
    if !has_rows {
        return "", errors.New("There are no rows in the result set")
    }
    result, scan_err := ScanAsString(stmnt)
    if scan_err != nil {
        return "", scan_err
    }
    return result[0], nil
}

func (c *Conn) ExecFirstAsInt( sql string ) (int, error) {
    stmnt, prep_err := c.Prepare(sql)
    if prep_err != nil {
        return 0, prep_err
    }
    has_rows := stmnt.Next()
    if !has_rows {
        return 0, errors.New("There are no rows in the result set")
    }
    result, scan_err := ScanAsInt(stmnt)
    if scan_err != nil {
        return 0, scan_err
    }
    return result[0], nil
}

