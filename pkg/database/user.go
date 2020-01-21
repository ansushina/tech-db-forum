package database

import (
	"errors"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/jackc/pgx"
)

var (
	UserIsExist  = errors.New("User is already exist")
	UserNotFound = errors.New("User not found")
	UserConflict = errors.New("User conflict")
)

func CreateUser(user *models.User) (*models.Users, error) {

	err := DB.DBPool.QueryRow(
		`
			INSERT INTO users (nickname, fullname, about, email)
			VALUES ($1, $2, $3, $4) 
			RETURNING nickname
		`,
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	).Scan(&user.Nickname)

	if err != nil {
		users := models.Users{}
		queryRows, err := DB.DBPool.Query(`
				SELECT "nickname", "fullname", "email", "about"
				FROM users
				WHERE "nickname" = $1 OR "email" = $2
			`,
			&user.Nickname,
			&user.Email)
		defer queryRows.Close()

		if err != nil {
			return nil, err
		}

		for queryRows.Next() {
			user := models.User{}
			queryRows.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
			users = append(users, &user)
		}
		return &users, UserIsExist
	}
	return nil, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := DB.DBPool.QueryRow(
		`SELECT nickname, fullname, about, email 
		FROM users WHERE email = $1`,
		email,
	).Scan(
		&u.Nickname,
		&u.Fullname,
		&u.About,
		&u.Email,
	)
	if err != nil {
		return models.User{}, UserNotFound
	}

	return u, nil
}

func GetUserInfo(nickname string) (models.User, error) {
	var u models.User
	err := DB.DBPool.QueryRow(
		`SELECT nickname, fullname, about, email 
		FROM users WHERE nickname = $1`,
		nickname,
	).Scan(
		&u.Nickname,
		&u.Fullname,
		&u.About,
		&u.Email,
	)
	if err != nil {
		return models.User{}, UserNotFound
	}

	return u, nil
}
func UpdateUser(user models.User) (models.User, error) {
	err := DB.DBPool.QueryRow(
		`
			UPDATE users
			SET fullname = coalesce(nullif($2, ''), fullname),
				email = coalesce(nullif($3, ''), email),
				about = coalesce(nullif($4, ''), about)
			WHERE "nickname" = $1
			RETURNING nickname, fullname, email, about
		`,
		&user.Nickname,
		&user.Fullname,
		&user.Email,
		&user.About,
	).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.Email,
		&user.About,
	)

	if err != nil {
		if ErrorCode(err) != pgxOK {
			return models.User{}, UserConflict
		}
		return models.User{}, UserNotFound
	}

	return user, nil
}

var queryForumUserWithSience = map[string]string{
	"true": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			) 
			AND LOWER(nickname) < LOWER($2::TEXT) COLLATE "POSIX"
		ORDER BY nickname  COLLATE "POSIX" DESC
		LIMIT $3::TEXT::INTEGER
	`,
	"false": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			) 
			AND LOWER(nickname) > LOWER($2::TEXT) COLLATE "POSIX"
		ORDER BY nickname  COLLATE "POSIX"
		LIMIT $3::TEXT::INTEGER
	`,
}

var queryForumUserNoSience = map[string]string{
	"true": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			)
		ORDER BY nickname  COLLATE "POSIX" DESC
		LIMIT $2::TEXT::INTEGER
	`,
	"false": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			)
		ORDER BY nickname  COLLATE "POSIX"
		LIMIT $2::TEXT::INTEGER
	`,
}

func GetForumUsers(slug, limit, since, desc string) (*models.Users, error) {
	var rows *pgx.Rows
	var err error

	if since != "" {
		rows, err = DB.DBPool.Query(queryForumUserWithSience[desc], slug, since, limit)
	} else {
		rows, err = DB.DBPool.Query(queryForumUserNoSience[desc], slug, limit)
	}
	defer rows.Close()

	if err != nil {
		return nil, ForumNotFound
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
			return nil, ForumNotFound
		}
	}
	return &users, nil
}
