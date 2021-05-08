package main

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	sql.Register("numbers", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.CreateModule("numbers", &numberModule{})
		},
	})
	db, err := sql.Open("numbers", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("create virtual table vals using numbers")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select idx, val from vals LIMIT 50;")
	if err != nil {
		log.Fatal("error executing query:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var idx, val int64
		if err := rows.Scan(&idx, &val); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("idx=%d, val=%d\n", idx, val)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
