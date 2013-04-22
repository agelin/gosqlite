// Copyright 2013 Richard B. Lyman. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sqlite

func ScanAllAsString( stmnt *Stmt ) [][]string {
    result := [][]string{}
    for stmnt.Next() {
        result = append(result, ScanAsString(stmnt))
    }
    return result
}

func ScanAsString( stmnt *Stmt ) []string {
    result := []string{}
    for i:=1; len(result) == 0; i++ {
        args := make([]string,i)
        arg_addrs := make([]interface{},i)
        for i := range args {
            arg_addrs[i] = &args[i]
        }
        scan_err := stmnt.Scan(arg_addrs...)
        if scan_err == nil {
            result = args
        }
    }
    return result
}
