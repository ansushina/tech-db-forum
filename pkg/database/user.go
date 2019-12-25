package database

import (
	"errors"
	"log"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/jackc/pgx"
)

var (
	UserIsExist  = errors.New("User is already exist")
	UserNotFound = errors.New("User not found")
)

func CreateUser(user models.User) (models.User, error) {
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

	switch ErrorCode(err) {
	case pgxOK:
		return user, nil
	case pgxErrUnique:
		f, _ := GetUserInfo(user.Nickname)
		return f, UserIsExist
	case pgxErrNotNull:
		return models.User{}, UserNotFound
	default:
		return models.User{}, err
	}
}
func GetUserInfo(nickname string) (models.User, error) {
	var u models.User
	err := DB.DBPool.QueryRow(
		`SELECT nickname, fullname, about, email 
		FROM users WHERE LOWER(nickname) = LOWER($1)`,
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
func UpdateUser(u models.User) (models.User, error) {
	queryString := "UPDATE users SET "
	if u.Email != "" {
		queryString += "email = '" + u.Email
		if u.About != "" || u.Fullname != "" {
			queryString += "', "
		}
	}
	if u.Fullname != "" {
		queryString += "fullname = '" + u.Fullname
		if u.About != "" {
			queryString += "', "
		}
	}
	if u.About != "" {
		queryString += "fullname = '" + u.Fullname + "'"
	}
	queryString += `WHERE LOWER(nickname) = LOWER('` + u.Nickname + "') RETURNING nickname"
	err := DB.DBPool.QueryRow(queryString).Scan(&u.Nickname)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func GetForumUsers(slug, limit, since string, desc bool) (models.Users, error) {
	queryString := " SELECT nickname, fullname, about, email FROM users u join forum_users f on u.nickname = f.forum_user "
	queryString += " where forum = '" + slug + "' "

	if since != "" {
		queryString += " AND u.id > " + since + " "
	}
	queryString += " order by id "
	if desc {
		queryString += " DESC "
	}

	if limit != "" {
		queryString += " limit " + limit
	}

	log.Print(queryString)
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
		//_, err := GetForumBySlug(slug)
		//if err != nil {
		//	return models.Users{}, ForumNotFound
		//}
	}
	return users, nil

}
