package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Created_at time.Time    `json:"created_at"`
	Updated_at sql.NullTime `json:"updated_at"`
}
