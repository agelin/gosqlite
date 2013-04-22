// Copyright 2013 Richard B. Lyman. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sqlite

/*
#include <sqlite3.h>
*/
import "C"

func Columns( stmnt *Stmt ) int {
    return int(C.sqlite3_column_count(stmnt.stmt))
}

func ScanAllAsString( stmnt *Stmt ) [][]string {
    result := [][]string{}
    for stmnt.Next() {
        result = append(result, ScanAsString(stmnt))
    }
    return result
}

func ScanAsString( stmnt *Stmt ) []string {
    result := make([]string,Columns(stmnt))
    addrs := make([]interface{},Columns(stmnt))
    for i := range result {
        addrs[i] = &result[i]
    }
    stmnt.Scan(addrs...) // Ignoring the scan error here... bad!
    return result
}

func (c *Conn) PrepareAndScanAllAsString( sql string ) ([][]string, error) {
    stmnt, prep_err := c.Prepare(sql)
    if prep_err != nil {
        return nil, prep_err
    }
    return ScanAllAsString(stmnt), nil
}
