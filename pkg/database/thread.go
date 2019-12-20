package dbhandlers

import (
	"errors"
	"github.com/ansushina/tech-db-forum/app/models"
)

var (
	ThreadIsExist          = errors.New("Forum is already exist")
	ThreadNotFound         = errors.New("Forum not found")
)

func GetThreadBySlug(string slug) (models.Thread, error) {
	var t models.Thread

	err := DB.DBPool.QueryRow(`SELECT id, votes, created, slug, title, author, forum, message FROM forums WHERE slug = $1`, slug).Scan(
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
