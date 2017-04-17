package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/yuuki/hworq/pkg/data"
)

const (
	DefaultDBUserName = "hworq"
	DefaultDBName     = "hworq"
)

var (
	Schemas = []string{
		"data/schema/queue.sql",
	}
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

func (db *DB) CreateSchema() error {
	for _, schema := range Schemas {
		sql, err := data.Asset(schema)
		if err != nil {
			return errors.Wrapf(err, "get schema error: %v", schema)
		}
		_, err = db.Exec(fmt.Sprintf("%s", sql))
		if err != nil {
			return errors.Wrapf(err, "exec schema error: %s", sql)
		}
	}
	return nil
}
