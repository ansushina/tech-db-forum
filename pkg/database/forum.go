package database

import (
	"errors"
	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/jackc/pgx"
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
	err := DB.DBPool.QueryRow(`SELECT slug, title, "user", posts, threads FROM forums WHERE slug = $1`, slug).Scan(
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
		thread, _ := GetThreadBySlug(t.Slug)
		return thread, ForumIsExist
	case pgxErrNotNull:
		return models.Thread{}, UserNotFound
	default:
		return models.Thread{}, err
	}
}

func GetForumThreads(slug, limit, since string, desc bool) (models.Threads, error) {

	queryString := " SELECT author, created, forum, id, message, slug, title, votes FROM threads "
	queryString += " where forum = " + slug + " "

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
		return models.Threads{}, err
	}

	threads := models.Threads{}

	for rows.Next() {
		t := models.Thread{}
		err = rows.Scan(
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.Id,
			&t.Message,
			&t.Slug,
			&t.Title,
			&t.Votes,
		)
		threads = append(threads, &t)
	}

	if len(threads) == 0 {
		_, err := GetForumBySlug(slug)
		if err != nil {
			return models.Threads{}, ForumNotFound
		}
	}
	return threads, nil
}

func GetForumUsers(slug, limit, since string, desc bool) (models.Users, error) {
	queryString := " SELECT nickname, fullname, about, email FROM users "
	queryString += " where forum = " + slug + " "

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
		return models.Users{}, err
	}

	users := models.Users{}

	for rows.Next() {
		u := models.User{}
		err = rows.Scan(
			&u.Nickname,
			&u.Fullname,
			&u.About,
			&u.Email,
		)
		users = append(users, &u)
	}

	if len(users) == 0 {
		_, err := GetForumBySlug(slug)
		if err != nil {
			return models.Users{}, ForumNotFound
		}
	}
	return users, nil

}
