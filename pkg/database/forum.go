package database

import (
	"errors"

	"github.com/ansushina/tech-db-forum/app/models"
)

var (
	ForumIsExist          = errors.New("Forum is already exist")
	ForumNotFound         = errors.New("Forum not found")
	ForumOrAuthorNotFound = errors.New("Forum or Author not found")
)

func CreateForum(forum models.Forum) (models.Forum, error) {
	err := DB.DBPool.QueryRow(
		`
			INSERT INTO forums (slug, title, "user")
			VALUES ($1, $2, $3) 
			RETURNING "user"
		`,
		&forum.Slug,
		&forum.Title,
		&forum.User,
	).Scan(&forum.User)

	switch ErrorCode(err) {
	case pgxOK:
		return forum, nil
	case pgxErrUnique:
		f, _ := GetForumBySlug(forum.Slug)
		return f, ForumIsExist
	case pgxErrNotNull:
		return models.Forum{}, UserNotFound
	default:
		return models.Forum{}, err
	}
}

func GetForumBySlug(slug string) (models.Forum, error) {
	var f models.Forum
	var p, t int

	err := DB.DBPool.QueryRow(`SELECT slug, title, "user", posts, threads FROM forums WHERE slug = $1`, slug).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&p,
		&t,
	)

	f.Posts = float32(p)
	f.Threads = float32(t)
	if err != nil {
		return models.Forum{}, err
	}

	return f, nil
}
