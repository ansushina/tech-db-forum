package dbhandlers

import (
	"errors"
	"github.com/ansushina/tech-db-forum/app/models"
)

var (
	ForumIsExist          = errors.New("Forum is already exist")
	ForumNotFound         = errors.New("Forum not found")
	ForumOrAuthorNotFound = errors.New("Forum or Author not found")
	UserNotFound          = errors.New("User not found")
)

func CreateForum(forum models.Forum) (models.Forum, error) {
	err := DB.DBPool.QueryRow(
		`
			INSERT INTO forums (slug, title, user)
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
	err := DB.DBPool.QueryRow(`SELECT * FROM forums WHERE slug = $1`, slug).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&f.Posts,
		&f.Threads,
	)
	if err != nil {
		return models.Forum{}, ForumNotFound
	}

	return f, nil
}

func CreateForumThread(t models.Thread) (models.Thread, error) {
	err := DB.DBPool.QueryRow(
		`INSERT INTO threads (title, author, forum, message, slug) 
		values ($1, $2, $3, $4, $5)`,
		&t.Title,
		&t.Author,
		&t.Forum,
		&t.Message,
		&t.Slug,
	).Scan(
		&t.Id,
		&t.Title,
		&t.Author,
		&t.Forum,
		&t.Message,
		&t.Slug,
		&t.Created,
		&t.Votes,
	)

	switch ErrorCode(err) {
	case pgxOK:
		return t, nil
	case pgxErrUnique:
		thread, _ := GetTreadBySlug(t.Slug)
		return thread, ForumIsExist
	case pgxErrNotNull:
		return models.Thread{}, UserNotFound
	default:
		return models.Thread{}, err
	}
}

func GetForumThreads() {

}

func GetForumUsers() {

}
