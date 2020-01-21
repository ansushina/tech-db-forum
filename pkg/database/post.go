package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/jackc/pgx"
)

var (
	PostIsExist    = errors.New("Post is already exist")
	PostNotFound   = errors.New("Post not found")
	ParentNotExist = errors.New("post parent not exist")
)

func GetPostById(id int) (models.Post, error) {
	var p models.Post
	err := DB.DBPool.QueryRow(
		`SELECT id, author, created, forum, "isEdited", message, parent, thread 
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

func GetPostFull(id int, related []string) (*models.PostFull, error) {
	postFull := models.PostFull{}
	var err error
	postFull.Post = &models.Post{}

	*postFull.Post, err = GetPostById(id)
	if err != nil {
		return nil, err
	}

	for _, obj := range related {
		switch obj {
		case "user":
			postFull.Author = &models.User{}
			*postFull.Author, err = GetUserInfo(postFull.Post.Author)
		case "forum":
			postFull.Forum = &models.Forum{}
			*postFull.Forum, err = GetForumBySlug(postFull.Post.Forum)
		case "thread":
			postFull.Thread = &models.Thread{}
			*postFull.Thread, err = GetThreadBySlug(strconv.Itoa(postFull.Post.Thread))
		}

		if err != nil {
			return nil, err
		}
	}

	return &postFull, nil
}

func UpdatePost(p models.Post) (models.Post, error) {
	post, e := GetPostById(p.Id)
	if e != nil {
		return models.Post{}, e
	}

	if len(p.Message) == 0 {
		return post, nil
	}
	err := DB.DBPool.QueryRow(`UPDATE posts SET message = COALESCE($1, message),
	"isEdited" = ($1 IS NOT NULL AND $1 <> message) WHERE id = $2
	RETURNING id, parent, author, message, "isEdited", forum, thread, created`, &p.Message, p.Id).Scan(&post.Id,
		&post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return models.Post{}, PostNotFound
	}
	p.IsEdited = true
	return post, nil
}

func authorExists(nickname string) bool {
	var user models.User
	err := DB.DBPool.QueryRow(
		`
			SELECT "nickname", "fullname", "email", "about"
			FROM users
			WHERE "nickname" = $1
		`,
		nickname,
	).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)

	if err != nil {
		return true
	}
	return false
}

const postID = `
	SELECT id
	FROM posts
	WHERE id = $1 AND thread IN (SELECT id FROM threads WHERE thread <> $2)
`

func parentExitsInOtherThread(parent int, threadID int) bool {
	var t int64
	err := DB.DBPool.QueryRow(postID, parent, threadID).Scan(&t)

	if err != nil {
		return false
	}
	return true
}

func parentNotExists(parent int) bool {
	if parent == 0 {
		return false
	}

	var t int64
	err := DB.DBPool.QueryRow(`SELECT id FROM posts WHERE id = $1`, parent).Scan(&t)

	if err != nil {
		return true
	}
	return false
}

func checkPost(p *models.Post, t *models.Thread) error {
	if authorExists(p.Author) {
		return UserNotFound
	}
	if parentExitsInOtherThread(p.Parent, t.Id) || parentNotExists(p.Parent) {
		return ParentNotExist
	}
	return nil
}

func CreateThreadPost(posts *models.Posts, param string) (*models.Posts, error) {
	thread, err := GetThreadBySlug(param)
	if err != nil {
		return nil, err
	}

	postsNumber := len(*posts)
	if postsNumber == 0 {
		return posts, nil
	}

	dateTimeTemplate := "2006-01-02 15:04:05"
	created := time.Now().Format(dateTimeTemplate)
	query := strings.Builder{}
	query.WriteString("INSERT INTO posts (author, created, message, thread, parent, forum, path) VALUES ")
	queryBody := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (SELECT last_value FROM posts_id_seq)),"
	for i, post := range *posts {
		err = checkPost(&post, &thread)
		if err != nil {
			return nil, err
		}

		temp := fmt.Sprintf(queryBody, post.Author, created, post.Message, thread.Id, post.Parent, thread.Forum, post.Parent)
		if i == postsNumber-1 {
			temp = temp[:len(temp)-1]
		}
		query.WriteString(temp)
	}
	query.WriteString("RETURNING author, created, forum, id, message, parent, thread")
	tx, txErr := DB.DBPool.Begin()
	if txErr != nil {
		return nil, txErr
	}
	defer tx.Rollback()

	rows, err := tx.Query(query.String())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	insertPosts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		rows.Scan(
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.Id,
			&post.Message,
			&post.Parent,
			&post.Thread,
		)
		insertPosts = append(insertPosts, post)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tx.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(insertPosts), thread.Forum)
	for _, p := range insertPosts {
		tx.Exec(`INSERT INTO forum_users VALUES ($1, $2) ON CONFLICT DO NOTHING`, p.Author, p.Forum)
	}
	tx.Commit()

	return &insertPosts, nil
}

var queryPostsWithSience = map[string]map[string]string{
	"true": map[string]string{
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND (path < (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
			ORDER BY path DESC
			LIMIT $3::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts p
			WHERE p.thread = $1 and p.path[1] IN (
				SELECT p2.path[1]
				FROM posts p2
				WHERE p2.thread = $1 AND p2.parent = 0 and p2.path[1] < (SELECT p3.path[1] from posts p3 where p3.id = $2)
				ORDER BY p2.path DESC
				LIMIT $3
			)
			ORDER BY p.path[1] DESC, p.path[2:]
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND id < $2::TEXT::INTEGER
			ORDER BY id DESC
			LIMIT $3::TEXT::INTEGER
		`,
	},
	"false": map[string]string{
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND (path > (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
			ORDER BY path
			LIMIT $3::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts p
			WHERE p.thread = $1 and p.path[1] IN (
				SELECT p2.path[1]
				FROM posts p2
				WHERE p2.thread = $1 AND p2.parent = 0 and p2.path[1] > (SELECT p3.path[1] from posts p3 where p3.id = $2::TEXT::INTEGER)
				ORDER BY p2.path
				LIMIT $3::TEXT::INTEGER
			)
			ORDER BY p.path
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND id > $2::TEXT::INTEGER
			ORDER BY id
			LIMIT $3::TEXT::INTEGER
		`,
	},
}

var queryPostsNoSience = map[string]map[string]string{
	"true": map[string]string{
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 
			ORDER BY path DESC
			LIMIT $2::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND path[1] IN (
				SELECT path[1]
				FROM posts
				WHERE thread = $1
				GROUP BY path[1]
				ORDER BY path[1] DESC
				LIMIT $2::TEXT::INTEGER
			)
			ORDER BY path[1] DESC, path
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1
			ORDER BY id DESC
			LIMIT $2::TEXT::INTEGER
		`,
	},
	"false": map[string]string{
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 
			ORDER BY path
			LIMIT $2::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND path[1] IN (
				SELECT path[1] 
				FROM posts 
				WHERE thread = $1 
				GROUP BY path[1]
				ORDER BY path[1]
				LIMIT $2::TEXT::INTEGER
			)
			ORDER BY path
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 
			ORDER BY id
			LIMIT $2::TEXT::INTEGER
		`,
	},
}

func GetThreadPosts(param, limit, since, sort, desc string) (*models.Posts, error) {
	thread, err := GetThreadBySlug(param)
	if err != nil {
		return nil, ThreadNotFound
	}

	var rows *pgx.Rows

	if since != "" {
		query := queryPostsWithSience[desc][sort]
		rows, err = DB.DBPool.Query(query, thread.Id, since, limit)
	} else {
		query := queryPostsNoSience[desc][sort]
		rows, err = DB.DBPool.Query(query, thread.Id, limit)
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &posts, nil
}
