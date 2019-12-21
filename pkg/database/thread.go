package database

import (
	"errors"
	"github.com/ansushina/tech-db-forum/app/models"
)

var (
	ThreadIsExist  = errors.New("Thread is already exist")
	ThreadNotFound = errors.New("Thread not found")
)

func GetThreadBySlug(slug string) (models.Thread, error) {
	var t models.Thread

	err := DB.DBPool.QueryRow(`SELECT id, votes, created, slug, title, author, forum, message FROM threads WHERE slug = $1`, slug).Scan(
		&t.Id,
		&t.Votes,
		&t.Created,
		&t.Slug,
		&t.Title,
		&t.Author,
		&t.Forum,
		&t.Message,
	)
	if err != nil {
		return models.Thread{}, ThreadNotFound
	}

	return t, nil
}

func GetThreadById(id int) (models.Thread, error) {
	var t models.Thread

	err := DB.DBPool.QueryRow(`SELECT id, votes, created, slug, title, author, forum, message FROM threads WHERE id = $1`, id).Scan(
		&t.Id,
		&t.Votes,
		&t.Created,
		&t.Slug,
		&t.Title,
		&t.Author,
		&t.Forum,
		&t.Message,
	)
	if err != nil {
		return models.Thread{}, ThreadNotFound
	}

	return t, nil
}
