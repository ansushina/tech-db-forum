package database

import (
	"io/ioutil"
	"log"
	"strconv"

	"github.com/jackc/pgx"
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
		Database: "docker",
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
	log.Printf("Connection created")

	content, err := ioutil.ReadFile("./database/create.sql")
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
	log.Printf("tables created")

	return nil
}

func ErrorCode(err error) string {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return pgxOK
	}
	return pgerr.Code
}

func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}
