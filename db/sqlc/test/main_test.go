package db

import (
	"database/sql"
	"log"
	"os"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..")

	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to Postgres: ", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}
