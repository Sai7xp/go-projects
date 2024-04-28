package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jaswdr/faker/v2"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:sumanth123@localhost:5432/xbank?sslmode=disable"
)

var testQueries *Queries
var fake faker.Faker
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Error while opening new db connection ", err)
	}

	// defer conn.Close()

	fake = faker.New()
	testQueries = New(testDB)
	os.Exit(m.Run())
}
