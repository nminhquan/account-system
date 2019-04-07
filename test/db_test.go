package mas_test

import (
	"database/sql"
	"fmt"
	"testing"

	"gitlab.zalopay.vn/quannm4/mas/db"
)

type FakeDB struct {
	*db.MasDB
}

func TestDBConnection(t *testing.T) {
	dbTest := db.CreateMySQLDB("localhost", "root", "123456", "mas_test")
	if dbTest == nil {
		t.Error("create DB error")
	} else {
		//test
		version := RunTestQuery(dbTest.DB())
		if version == "" {
			t.Error("run query error")
		}
	}
}

func RunTestQuery(masDB *sql.DB) string {
	results, err := masDB.Query("SELECT 1")

	if err != nil {
		panic(err.Error())
	}
	var version string
	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&version)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		fmt.Printf("TEST QUERY MYSQL version: [%v]\n", version)
	}
	return version
}
