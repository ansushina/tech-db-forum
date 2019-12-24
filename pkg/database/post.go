package database

import (
	"errors"
	"github.com/jackc/pgx"

	"github.com/ansushina/tech-db-forum/app/models"
)

var (
	PostIsExist  = errors.New("Post is already exist")
	PostNotFound = errors.New("Post not found")
)

func GetPostById(id int) (models.Post, error) {
	var p models.Post
	err := DB.DBPool.QueryRow(
		`SELECT id, author, created, forum, isEdited, message, parent, thread 
		FROM posts WHERE id = $1`, id).Scan(
		&p.Id,
		&p.Author,
		&p.Created,
		&p.Forum,
		&p.IsEdited,
		&p.Message,
		&p.Parent,
		&p.Thread,
	)
	if err != nil {
		return models.Post{}, PostNotFound
	}

	return p, nil
}
func UpdatePost(p models.Post) (models.Post, error) {
	err := DB.DBPool.QueryRow(
		`UPDATE users SET message = $2, isEdited = 'true' 
		WHERE id = $1`,
		&p.Id,
		&p.Message,
	).Scan()
	if err != nil {
		return models.Post{}, UserNotFound
	}
	p.IsEdited = true
	return p, nil
}

func CreateThreadPost(slug string, p models.Post) (models.Post, error) {
	err := DB.DBPool.QueryRow(
		`
			INSERT INTO posts (author, forum, message, parent, thread)
			VALUES ($1, $2, $3) 
			RETURNING "id"
		`,
		&p.Author,
		&p.Forum,
		&p.Message,
		&p.Parent,
		&p.Thread,
	).Scan(&p.Id)

	switch ErrorCode(err) {
	case pgxOK:
		return p, nil
	case pgxErrNotNull:
		return models.Post{}, UserNotFound
	default:
		return models.Post{}, err
	}
}

func GetThreadPosts(param string, limit, since string, desc bool) (models.Posts, error) {
	queryString := " SELECT author, created, forum, id, message, parent, thread, isEdited FROM posts "
	if isNumber(param) {
		queryString += " where thread = " + param + " "
	} else {
		return models.Posts{}, errors.New("Wrong params")
	}

	if since != "" {
		queryString += " AND t.created <= TIMESTAMPTZ '" + since + "' "
	}
	queryString += " order by created "
	if desc {
		queryString += " DESC "
	}

	if limit != "" {
		queryString += " limit " + limit
	}

	var rows *pgx.Rows
	var err error
	rows, err = DB.DBPool.Query(queryString)

	if err != nil {
		return models.Posts{}, err
	}

	posts := models.Posts{}

	for rows.Next() {
		t := models.Post{}
		err = rows.Scan(
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.Id,
			&t.Message,
			&t.Parent,
			&t.Thread,
			&t.IsEdited,
		)
		posts = append(posts, &t)
	}

	if len(posts) == 0 {
		_, err := GetThreadBySlug(param)
		if err != nil {
			return models.Posts{}, ForumNotFound
		}
	}
	return posts, nil
}
