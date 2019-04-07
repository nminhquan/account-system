package mas_test

import (
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.zalopay.vn/quannm4/mas/db"
)

const TEST_KEY = "__key__"
const TEST_VALUE = 1.0

func TestRocksDBBDD(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RocksDB test")
}

var newDB = func() *db.RocksDB {
	db, _ := db.NewRocksDB("./test_db/abc")
	return db
}

var _ = Describe("RocksDB test", func() {
	var (
		testDB *db.RocksDB = newDB()
	)

	AfterSuite(func() {
		testDB.Close()
	})

	It("should set account successfully", func() {
		err := testDB.SetAccountBalance(TEST_KEY, TEST_VALUE)
		Expect(err).Should(BeNil())
	})
	It("should get account successfully", func() {
		s, _ := testDB.GetAccountBalance(TEST_KEY)
		number, err := strconv.ParseFloat(s, 64)
		Expect(err).Should(BeNil())
		Expect(number).Should(Equal(TEST_VALUE))
	})
	It("should be Key does not exist", func() {
		s, err := testDB.GetAccountBalance("some random key not exist")
		Expect(s).To(BeEmpty())
		Expect(err).ToNot(BeNil())
		Expect(err).To(Equal(db.ErrRockDBKeyNotExist))
	})
})
