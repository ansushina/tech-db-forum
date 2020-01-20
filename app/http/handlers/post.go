package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/ansushina/tech-db-forum/pkg/database"
)

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	id_str, _ := checkVar("id", r)
	id, _ := strconv.Atoi(id_str)

	query := r.URL.Query()
	related := strings.Split(query.Get("related"), ",")

	res, err := database.GetPostFull(id, related)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.PostNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find post with id "+id_str)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {

	p := models.Post{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = p.UnmarshalJSON(body)

	id, _ := checkVar("id", r)
	id_f, _ := strconv.Atoi(id)
	p.Id = int(id_f)

	res, err := database.UpdatePost(p)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.PostNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find post with id "+id)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func PostsCreate(w http.ResponseWriter, r *http.Request) {
	slug, _ := checkVar("slug", r)

	p := models.Posts{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = p.UnmarshalJSON(body)

	_, err := database.GetThreadBySlug(slug)

	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
		return
	}

	if len(p) == 0 {
		l := models.Threads{}
		WriteResponse(w, http.StatusCreated, l)
		return
	}

	_, err = database.GetUserInfo(p[0].Author)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+p[0].Author)
		return
	}

	res, e := database.CreateThreadPost(&p, slug)
	if e == database.ParentNotExist {
		WriteErrorResponse(w, http.StatusConflict, e.Error())
		return
	} else if e != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, e.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, *res)
}
