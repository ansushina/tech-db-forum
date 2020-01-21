package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/ansushina/tech-db-forum/pkg/database"
)

func ForumCreate(w http.ResponseWriter, r *http.Request) {

	f := models.Forum{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = f.UnmarshalJSON(body)

	newForum, err := database.CreateForum(f)
	if err == database.UserNotFound {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+f.User)
		return
	} else if err == database.ForumIsExist {
		WriteResponse(w, http.StatusConflict, newForum)
		return
	} else if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, newForum)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	slug, _ := checkVar("slug", r)

	f, err := database.GetForumBySlug(slug)
	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, f)
	case database.ForumNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find forum with slug"+slug)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	slug, _ := checkVar("slug", r)
	query := r.URL.Query()
	since := query.Get("since")
	limit := query.Get("limit")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	res, err := database.GetForumThreads(slug, limit, since, desc)

	if err == database.ForumNotFound {
		WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusOK, res)

}

func ForumGetUsers(w http.ResponseWriter, r *http.Request) {
	slug, _ := checkVar("slug", r)
	query := r.URL.Query()
	limit := query.Get("limit")
	if limit == "" {
		limit = "100"
	}
	since := query.Get("since")
	desc := query.Get("desc")
	if desc == "" {
		desc = "false"
	}

	res, err := database.GetForumUsers(slug, limit, since, desc)

	if err == database.ForumNotFound {
		WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, res)
}
