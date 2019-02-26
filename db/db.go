package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func CreateDB(host string, userName string, password string) *MasDB {
	// var dsn string = fmt.Sprintf("%s:%s@/")
	db, error := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/mas", userName, password, host))
	// if there is an error opening the connection, handle it
	if error != nil {
		panic(error.Error())
	}
	masDB := MasDB{db}
	return &masDB
}
