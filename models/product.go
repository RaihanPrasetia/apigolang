package models

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Price       int          `json:"price"`
	User_id     int          `json:"user_id"`
	Category_id int          `json:"category_id"`
	Created_at  time.Time    `json:"created_at"`
	Updated_at  sql.NullTime `json:"updated_at"`
}
