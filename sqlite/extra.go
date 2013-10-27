// Copyright 2013 Richard B. Lyman. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sqlite

/*
#include <sqlite3.h>
#include <stdlib.h>



*/
import "C"
import "errors"
import "fmt"
import "log"
import "unsafe"

func ScanAsString(stmnt *Stmt) ([]string, error) {
	result := make([]string, Columns(stmnt))
    for i := range result {
        p := C.sqlite3_column_text(stmnt.stmt, C.int(i))
        n := C.sqlite3_column_bytes(stmnt.stmt, C.int(i))
        temp := C.GoStringN((*C.char)(unsafe.Pointer(p)), n)
        if len(temp) == 0 {
        }
        result[i] = temp
        //result[i] = string(make([]byte, n))
        //result[i] = "bob"
    }
	return result, nil
}

// ScanAllAsString will return all rows if there are no errors.
// If there are errors the rows successfully scanned up to
// that point are the only rows returned.
func ScanAllAsString(stmnt *Stmt) ([][]string, error) {
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

// See ScanAllAsString for an explanation of how possible errors
// affect the rows returned
func (c *Conn) ExecToStrings(sql string) ([][]string, error) {
	stmnt, prep_err := c.Prepare(sql)
	if prep_err != nil {
		return nil, prep_err
	}
	defer stmnt.Finalize()
	has_rows := stmnt.Next()
	if !has_rows {
		return [][]string{}, stmnt.Error()
	}
	stmnt.Reset()
	return ScanAllAsString(stmnt)
}

func (c *Conn) Throwaway(sql string) {
	stmnt, err := c.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmnt.Finalize()
	if err = stmnt.Exec(); err != nil {
		panic(err)
	}
	for stmnt.Next() {
	}
}
/*
func (c *Conn) RestrictedDump() ([]byte, error) {
//	tableSql, tableSql_err := c.ExecToStrings("SELECT sql FROM sqlite_master WHERE type='table'")
//    if tableSql_err != nil {
//        return nil, tableSql_err
//    }
    // Reduce each row's single value to a series of strings
	tables, tables_err := c.ExecToStrings("SELECT name FROM sqlite_master WHERE type='table'")
    if tables_err != nil {
        return nil, tables_err
    }
    results := []byte{}
    for _, tableRow := range tables {
        table := tableRow[0]
        fmt.Println("Gathering data from table:", table)
        results = append( results, "INSERT INTO " + table)
        // Get column names... from select * limit 1... then result[C.GoString(C.sqlite3_column_name(stmnt.stmt, C.int(i)))] = v
        rows, rows_err := c.ExecToStringMaps("SELECT * FROM "+table)
        if rows_err != nil {
            return nil, rows_err
        }
        for _, row := range rows {
            fmt.Println("\trow:", row)
            for k, v := range row {
                fmt.Println("\t\t", k, ":", v)
            }
        }
        results = append(results, ";")
    }
    return nil, nil

//'insert into ' + t + ' (' + column_names_in_csv + ') values (' + values_in_csv + ')'
*/

/*
	stmnt, prep_err := c.Prepare(sql)
	if prep_err != nil {
		return nil, prep_err
	}
	defer stmnt.Finalize()
	has_rows := stmnt.Next()
	if !has_rows {
		return [][]string{}, stmnt.Error()
	}
	stmnt.Reset()
	return ScanAllAsString(stmnt)
*/

//}

func (c *Conn) DropAllTables() {
	c.Throwaway("PRAGMA writable_schema = 1")
	c.Throwaway("DELETE FROM sqlite_master WHERE type='table'")
	c.Throwaway("PRAGMA writable_schema = 0")
	c.Throwaway("VACUUM")
	c.Throwaway("PRAGMA integrity_check")
}

func Columns(stmnt *Stmt) int {
	return int(C.sqlite3_column_count(stmnt.stmt))
}

func ScanAllAsStringMap(stmnt *Stmt) ([]map[string]string, error) {
	result := []map[string]string{}
	for stmnt.Next() {
		row, err := ScanAsMap(stmnt)
		if err != nil {
			return result, err
		}
		result = append(result, row)
	}
	return result, nil
}

func ScanAsMap(stmnt *Stmt) (map[string]string, error) {
	temp_result := make([]string, Columns(stmnt))
	addrs := make([]interface{}, Columns(stmnt))
	for i := range temp_result {
		addrs[i] = &temp_result[i]
	}
	err := stmnt.Scan(addrs...)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	for i, v := range temp_result {
		result[C.GoString(C.sqlite3_column_name(stmnt.stmt, C.int(i)))] = v
	}
	return result, nil
}

func ScanAsInt(stmnt *Stmt) ([]int, error) {
	result := make([]int, Columns(stmnt))
	addrs := make([]interface{}, Columns(stmnt))
	for i := range result {
		addrs[i] = &result[i]
	}
	err := stmnt.Scan(addrs...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Conn) ExecToStringMaps(sql string) ([]map[string]string, error) {
	stmnt, prep_err := c.Prepare(sql)
	if prep_err != nil {
		return nil, prep_err
	}
	defer stmnt.Finalize()
	has_rows := stmnt.Next()
	if !has_rows {
		return []map[string]string{}, stmnt.Error()
	}
	stmnt.Reset()
	return ScanAllAsStringMap(stmnt)
}

func (c *Conn) SafeExecToStrings(sql string) ([][]string, error) {
	begin_err := c.Exec("BEGIN;")
	if begin_err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to begin safe execution: %s", begin_err))
	}
	result, result_err := c.ExecToStrings(sql)
	if result_err != nil {
		early_rollback_err := c.Exec("ROLLBACK;")
		if early_rollback_err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to rollback after failing to execute query: %s, %s", result_err, early_rollback_err))
		}
		return nil, errors.New(fmt.Sprintf("Failed to execute query: %s", result_err))
	}
	rollback_err := c.Exec("ROLLBACK;")
	if rollback_err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to rollback safe execution: %s", rollback_err))
	}
	return result, nil
}

func (c *Conn) SafeExecToStringMaps(sql string) ([]map[string]string, error) {
	begin_err := c.Exec("BEGIN;")
	if begin_err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to begin safe execution: %s", begin_err))
	}
	result, result_err := c.ExecToStringMaps(sql)
	if result_err != nil {
		early_rollback_err := c.Exec("ROLLBACK;")
		if early_rollback_err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to rollback after failing to execute query: %s, %s", result_err, early_rollback_err))
		}
		return nil, errors.New(fmt.Sprintf("Failed to execute query: %s", result_err))
	}
	rollback_err := c.Exec("ROLLBACK;")
	if rollback_err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to rollback safe execution: %s", rollback_err))
	}
	return result, nil
}

func (c *Conn) ExecFirstAsString(sql string) (string, error) {
	stmnt, prep_err := c.Prepare(sql)
	if prep_err != nil {
		return "", prep_err
	}
	defer stmnt.Finalize()
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

func (c *Conn) ExecFirstAsInt(sql string) (int, error) {
	stmnt, prep_err := c.Prepare(sql)
	if prep_err != nil {
		return 0, prep_err
	}
	defer stmnt.Finalize()
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
