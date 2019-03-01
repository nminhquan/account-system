package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type MasDB struct {
	db *sql.DB
}

func (masDB *MasDB) DB() *sql.DB {
	return masDB.db
}

func (masDB *MasDB) Close() {
	masDB.db.Close()
}

func CreateMySQLDB(host string, userName string, password string, dbName string) *MasDB {
	// var dsn string = fmt.Sprintf("%s:%s@/")
	db, error := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", userName, password, host, dbName))
	// if there is an error opening the connection, handle it
	if error != nil {
		panic(error.Error())
	}
	masDB := MasDB{db}
	return &masDB
}

func CreateSQLite3DB(host string, userName string, password string, dbName string) *MasDB {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s.db", dbName))
	if err != nil {
		panic(err)
	}
	masDB := MasDB{db}
	return &masDB
}
