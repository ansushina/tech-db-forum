package database

import (
	"errors"
	"github.com/ansushina/tech-db-forum/app/models"
)

func ClearAll() error {

	res, err := DB.DBPool.Query("TRUNCATE TABLE forum_users, forum_forum, forum_thread, forum_post, forum_vote CASCADE")
	if err != nil {
		return errors.New("can't truncate tables")
	}
	defer res.Close()

	return nil
}

func DatabaseStatus() (models.Status, error) {
	res, err := DB.DBPool.Query("SELECT * FROM (SELECT count(posts) FROM forum_forum) as ff" +
		" CROSS JOIN (SELECT count(id) FROM forum_post) as fp" +
		" CROSS JOIN (SELECT count(id) FROM forum_thread) as ft" +
		" CROSS JOIN (SELECT count(nickname) FROM forum_users) as fu")
	defer res.Close()

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
