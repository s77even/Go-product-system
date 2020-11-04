package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//NewMysqlConn 创建mysql连接
func NewMysqlConn() (db *sql.DB, err error) {
	dsn := "root:seven7777777@tcp(127.0.0.1:3306)/imooc?charset=utf8"
	db, err = sql.Open("mysql", dsn)
	return
}

//GetResultRow 从数据库中获取一行
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				record[columns[i]] = string(v.([]byte))
			}
		}
	}
	return record
}

//GetResultRow 从数据库中获取多行
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	columns, _ := rows.Columns()
	vals := make([][]byte, len(columns))
	scans := make([]interface{}, len(columns))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := columns[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	return result
}
