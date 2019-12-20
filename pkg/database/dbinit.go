package dbhandlers

import (
	"io/ioutil"

	"github.com/jackc/pgx"
)

type DataBase struct {
	DBPool *pgx.ConnPool
}

var DB DataBase

func (db *DataBase) createConn() (err error) {
	conConfig := pgx.ConnConfig{
		Database: "zxc",
		User:     "docker",
		Password: "docker",
	}
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     conConfig,
		MaxConnections: 25,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	con, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return err
	}
	db.DBPool = con

	content, err := ioutil.ReadFile("./p")
	if err != nil {
		return err
	}

	tx, err := db.DBPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(string(content)); err != nil {
		return err
	}
	tx.Commit()
	return
}

func (db *DataBase) Close() error {
	if db.DBPool == nil {
		return errors.new("no connection")
	}
	db.DBPool.Close()
	return nil
}

func (db *DataBase) Begin() (tx *pgx.Tx, err error) {
	if db.DBPool == nil {
		return tx, errors.new("no connection")
	}

	return db.DBPool.Begin()
}

func (db *DataBase) QueryRow(query string, args ...interface{}) (row *pgx.Row, err error) {
	tx, err := Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row = tx.QueryRow(query, args...)

	return row, tx.Commit()
}

func (db *DataBase) Query(query string, args ...interface{}) (rows *pgx.Rows, err error) {
	tx, err := Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err = tx.Query(query, args...)
	if err != nil {
		return
	}

	return rows, tx.Commit()
}

func (db *DataBase) Exec(query string, args ...interface{}) (tag pgx.CommandTag, err error) {
	tx, err := Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	tag, err = tx.Exec(query, args...)
	if err != nil {
		return
	}

	return tag, tx.Commit()
}
