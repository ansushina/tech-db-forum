package database

import (
	"errors"
	"fmt"

	"github.com/ansushina/tech-db-forum/app/models"
)

var (
	ForumIsExist          = errors.New("Forum is already exist")
	ForumNotFound         = errors.New("Forum not found")
	ForumOrAuthorNotFound = errors.New("Forum or Author not found")
)

func CreateForum(forum models.Forum) (models.Forum, error) {
	var nick string
	err := DB.DBPool.QueryRow(
		`
			INSERT INTO forums (slug, title, "user")
			VALUES ($1, $2,  (
				SELECT nickname FROM users WHERE LOWER(nickname) = LOWER($3)
				))
			RETURNING "user"
		`,
		forum.Slug,
		forum.Title,
		forum.User,
	).Scan(&nick)
	forum.User = nick
	fmt.Println(nick)

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

	err := DB.DBPool.QueryRow(`SELECT slug, title, "user", posts, threads FROM forums WHERE LOWER(slug) = LOWER($1)`, slug).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&p,
		&t,
	)

	f.Posts = p
	f.Threads = t
	if err != nil {
		return models.Forum{}, ForumNotFound
	}

	return f, nil
}
