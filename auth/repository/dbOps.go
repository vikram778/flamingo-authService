package repository

import "github.com/jmoiron/sqlx"

type DBOps struct {
	db *sqlx.DB
}

func NewDBOpsRepository(db *sqlx.DB) *DBOps {
	return &DBOps{db: db}
}
