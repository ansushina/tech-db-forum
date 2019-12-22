package database

import (
	"errors"
	"github.com/jackc/pgx"
	"io/ioutil"
)

const (
	pgxOK            = ""
	pgxErrNotNull    = "23502"
	pgxErrForeignKey = "23503" // может возникнуть при добавлении дубликата
	pgxErrUnique     = "23505"
	pgxnoRows        = "no rows in result set"
)

type DataBase struct {
	DBPool *pgx.ConnPool
}

var DB DataBase

func CreateConn() (err error) {
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
	DB.DBPool = con

	content, err := ioutil.ReadFile("./p")
	if err != nil {
		return err
	}

	tx, err := DB.DBPool.Begin()
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
		return errors.New("no connection")
	}
	db.DBPool.Close()
	return nil
}

func Begin() (tx *pgx.Tx, err error) {
	if DB.DBPool == nil {
		return tx, errors.New("no connection")
	}

	return DB.DBPool.Begin()
}

func QueryRow(query string, args ...interface{}) (row *pgx.Row, err error) {
	tx, err := Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row = tx.QueryRow(query, args...)

	return row, tx.Commit()
}

func Query(query string, args ...interface{}) (rows *pgx.Rows, err error) {
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

func Exec(query string, args ...interface{}) (tag pgx.CommandTag, err error) {
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

func ErrorCode(err error) string {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return pgxOK
	}
	return pgerr.Code
}
