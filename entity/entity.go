package entity

import (
	"os"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"

	"gondos/jetgen/gondos/table"
)

// dbName used to create new table from generated jetgen.
// jet uses db name that is used during generation, but
// in runtime we may want to change db name.
var dbName string

func init() {
	c, err := mysqldriver.ParseDSN(os.Getenv("DSN"))
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("unable parse dbname from dsn")
	}
	dbName = c.DBName
	// change schema
	table.UseSchema(dbName)
}
