//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Users struct {
	ID              int64 `sql:"primary_key"`
	Name            string
	Email           *string
	CryptedPassword *string
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
