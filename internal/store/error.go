package store

import (
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/go-sql-driver/mysql"

	"gondos/internal/app"
)

// handleErr maps error to app Level errors
func (*UserStorage) handleErr(err error) error {
	var mysqlErr *mysql.MySQLError

	switch {
	case errors.Is(err, qrm.ErrNoRows):
		return fmt.Errorf("%w: %w", app.ErrUserNotFound, err)
	case errors.As(err, &mysqlErr):
		if mysqlErr.Number == 1062 { // DUPLICATE
			return app.ErrUserRegistered
		}
	}

	return err
}
