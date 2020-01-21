package database

import (
	"errors"
	"strconv"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/jackc/pgx"
)

var (
	ThreadIsExist  = errors.New("Thread is already exist")
	ThreadNotFound = errors.New("Thread not found")
)

func GetThreadBySlug(param string) (models.Thread, error) {
	if isNumber(param) {
		id := param
		var t models.Thread

		err := DB.DBPool.QueryRow(`
		SELECT id, votes, created, slug, title, author, forum, message 
		FROM threads WHERE id = $1`, id).Scan(
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
	} else {
		slug := param
		var t models.Thread

		err := DB.DBPool.QueryRow(`SELECT id, votes, created, slug, title, author, forum, message FROM threads WHERE LOWER(slug) = LOWER($1)`, slug).Scan(
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
}

func UpdateThreadBySlugorId(param string, t models.Thread) (models.Thread, error) {
	threadFound, err := GetThreadBySlug(param)
	if err != nil {
		return models.Thread{}, PostNotFound
	}

	updatedThread := models.Thread{}

	err = DB.DBPool.QueryRow(`
			UPDATE threads
			SET title = coalesce(nullif($2, ''), title),
				message = coalesce(nullif($3, ''), message)
			WHERE slug = $1
			RETURNING id, title, author, forum, message, votes, slug, created
		`,
		&threadFound.Slug,
		&t.Title,
		&t.Message,
	).Scan(
		&updatedThread.Id,
		&updatedThread.Title,
		&updatedThread.Author,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.Votes,
		&updatedThread.Slug,
		&updatedThread.Created,
	)

	if err != nil {
		return models.Thread{}, err
	}

	return updatedThread, nil
}

func VoteForThread(param string, vote models.Vote) (models.Thread, error) {
	var err error

	tx, txErr := DB.DBPool.Begin()
	if txErr != nil {
		return models.Thread{}, txErr
	}
	defer tx.Rollback()

	var thread models.Thread
	if isNumber(param) {
		id, _ := strconv.Atoi(param)
		err = tx.QueryRow(`SELECT id, author, created, forum, message, slug, title, votes FROM threads WHERE id = $1`, id).Scan(
			&thread.Id,
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
	} else {
		err = tx.QueryRow(`SELECT id, author, created, forum, message, slug, title, votes FROM threads WHERE LOWER(slug) = LOWER($1)`, param).Scan(
			&thread.Id,
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
	}
	if err != nil {
		return models.Thread{}, ForumNotFound
	}

	var nick string
	err = tx.QueryRow(`SELECT nickname FROM users WHERE nickname = $1`, vote.Nickname).Scan(&nick)
	if err != nil {
		return models.Thread{}, UserNotFound
	}

	rows, err := tx.Exec(`UPDATE votes SET voice = $1 WHERE thread = $2 AND nickname = $3;`,
		vote.Voice,
		thread.Id,
		vote.Nickname)
	if rows.RowsAffected() == 0 {
		_, err := tx.Exec(`INSERT INTO votes (nickname, thread, voice) VALUES ($1, $2, $3);`,
			vote.Nickname,
			thread.Id,
			vote.Voice)
		if err != nil {
			return models.Thread{}, UserNotFound
		}
	}

	err = tx.QueryRow(`SELECT votes FROM threads WHERE id = $1`, thread.Id).Scan(&thread.Votes)
	if err != nil {
		return models.Thread{}, err
	}

	tx.Commit()

	return thread, nil
}

func GetForumThreads(slug, limit, since string, desc bool) (models.Threads, error) {

	queryString := " SELECT author, created, forum, id, message, slug, title, votes FROM threads t "
	queryString += " where LOWER(forum) = LOWER('" + slug + "')"

	if since != "" && desc {
		queryString += " AND t.created <= TIMESTAMPTZ '" + since + "' "
	} else if since != "" {
		queryString += " AND t.created >= TIMESTAMPTZ '" + since + "' "
	}
	queryString += " order by created "
	if desc {
		queryString += " DESC "
	}

	if limit != "" {
		queryString += " limit " + limit
	}

	//log.Print(queryString)
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

	return threads, nil
}

func CreateForumThread(t models.Thread) (models.Thread, error) {
	err := DB.DBPool.QueryRow(
		`INSERT INTO threads (title, author, forum, message, slug, created) 
		values ($1, (select nickname from users where LOWER(nickname) = LOWER($2)), 
			(select slug from forums where LOWER(slug)=LOWER($3)), $4, $5, $6)
		RETURNING id, title, author, forum, message, slug, created, votes
		`,
		t.Title,
		t.Author,
		t.Forum,
		t.Message,
		t.Slug,
		t.Created,
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
		return thread, ThreadIsExist
	case pgxErrNotNull:
		return models.Thread{}, UserNotFound
	default:
		return models.Thread{}, err
	}
}
