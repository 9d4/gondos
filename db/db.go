package db

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	mysqldrv "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	return db, err
}

func NewLogWriter(file string) (io.Writer, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("db: unable create log file: %v", err)
	}

	writer := diode.NewWriter(f, 1024, 100*time.Millisecond, func(missed int) {
		log.Warn().Caller().
			Str("file", file).
			Int("missed", missed).
			Msg("logger dropped messages")
	})
	return writer, nil
}

func SetQueryLogger(output io.Writer) {
	log := log.Output(output)
	mysql.SetQueryLogger(func(ctx context.Context, info mysql.QueryInfo) {
		f, l, _ := info.Caller()
		log.Info().
			AnErr("err", info.Err).
			Dur("dur", info.Duration).
			Int64("rows", info.RowsProcessed).
			Str("file", f+":"+strconv.Itoa(l)).
			Str("query", info.Statement.DebugSql()).
			Send()
	})
}

func ParseSchema(dsn string) (string, error) {
	cfg, err := mysqldrv.ParseDSN(os.Getenv("DSN"))
	if err != nil {
		return "", err
	}
	return cfg.DBName, nil
}
