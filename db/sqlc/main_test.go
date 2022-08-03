package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/1BarCode/go-bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// var err error -> not needed anymore because err is defined on line 24 and gets overwritten by 29
	
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}