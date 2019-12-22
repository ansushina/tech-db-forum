package database

import (
	"errors"

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
