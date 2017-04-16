package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	DefaultDBUserName = "hworq"
	DefaultDBName     = "hworq"
)

type DB struct {
	*sql.DB
}

func New() (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable", DefaultDBUserName, DefaultDBName,
	))
	if err != nil {
		return nil, errors.Wrap(err, "postgres open error")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "postgres ping error")
	}
	return &DB{db}, nil
}
