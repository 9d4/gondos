package store

import (
	"errors"
	"gondos/internal/app"

	"github.com/go-sql-driver/mysql"
)

// handleErr maps error to app Level errors
func (s *UserStorage) handleErr(err error) error {
	var mysqlErr *mysql.MySQLError

	switch {
	case errors.As(err, &mysqlErr):
		if mysqlErr.Number == 1062 { // DUPLICATE
			return app.ErrUserRegistered
		}
	}

	return err
}
