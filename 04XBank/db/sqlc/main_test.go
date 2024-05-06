package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jaswdr/faker/v2"
	_ "github.com/lib/pq"
	"github.com/sai7xp/xbank/utils"
)

var testQueries *Queries
var fake faker.Faker
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatal("Can't load config: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error while opening new db connection ", err)
	}

	// defer conn.Close()

	fake = faker.New()
	testQueries = New(testDB)
	os.Exit(m.Run())
}
