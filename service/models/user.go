package models

import "time"

// User - is a Postgres user
type User struct {
	tableName struct{} `pg:"users"`

	UUID string `json:"uuid" pg:"uuid,notnull,pk"`
	Name string `json:"name" pg:"firstname,notnull"`
	Age  int    `json:"age" pg:"age,notnull"`

	CreatedAt time.Time `json:"-" pg:"created_at,notnull"`
	UpdatedAt time.Time `json:"-" pg:"updated_at"`
	DeletedAt time.Time `json:"-" pg:"deleted_at"`
}
