package mytest

import "testing"
import "database/sql"
import "mas/db"
import "fmt"

type FakeDB struct {
}

func createFakeDB() *FakeDB {
	fakeDB := FakeDB{}
	return &fakeDB
}
func TestDBConnection(t *testing.T) {
	dbTest := db.CreateDB("localhost", "root", "123456")
	if dbTest == nil {
		t.Error("create DB error")
	} else {
		//test
		version := RunQuery(dbTest.DB())
		if version == "" {
			t.Error("run query error")
		}
	}
}

func RunQuery(masDB *sql.DB) string {
	results, err := masDB.Query("SELECT @@VERSION")

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
