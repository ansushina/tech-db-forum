package dbhandlers

import (
	"github.com/jackc/pgx"
	"github.com/ansushina/tech-db-forum/app/models"
)


func CreateForum(forum *models.forum) (error) {
	err := DB.DBPool.QueryRow(
		`
			INSERT INTO forums (slug, title, user)
			VALUES ($1, $2, $3) 
			RETURNING "user"
		`,
		&f.Slug,
		&f.Title,
		&f.User,
	).Scan(&f.User)

	switch ErrorCode(err) {
		case pgxOK:
			return f, nil
		case pgxErrUnique:
			forum, _ := GetForum(f.Slug)
			return forum, ForumIsExist
		case pgxErrNotNull:
			return nil, UserNotFound
		default:
			return nil, err
	}
}
