package database

import (
	"errors"
	"log"

	"github.com/ansushina/tech-db-forum/app/models"
)

func ClearAll() error {

	res, err := DB.DBPool.Query("TRUNCATE TABLE users, forum_users, forums, threads, posts, votes CASCADE")
	if err != nil {
		return errors.New("can't truncate tables")
	}
	defer res.Close()

	return nil
}

func DatabaseStatus() (models.Status, error) {
	res, err := DB.DBPool.Query("SELECT * FROM (SELECT count(title) FROM forums) as ff" +
		" CROSS JOIN (SELECT count(id) FROM posts) as fp" +
		" CROSS JOIN (SELECT count(id) FROM threads) as ft" +
		" CROSS JOIN (SELECT count(nickname) FROM users) as fu")
	defer res.Close()

	log.Print(err)

	if err != nil {
		return models.Status{}, errors.New("cant get db statisics")
	}

	s := models.Status{}
	for res.Next() {
		err = res.Scan(&s.Forum, &s.Post, &s.Thread, &s.User)

		if err != nil {
			return models.Status{}, errors.New("db query result parsing error")
		}

		return s, nil
	}

	return models.Status{}, errors.New("cant get db statisics")
}
