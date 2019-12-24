package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/ansushina/tech-db-forum/pkg/database"
)

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	id_str, _ := checkVar("id", r)
	id, _ := strconv.Atoi(id_str)

	res, err := database.GetPostById(id)

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
	p.Id = float32(id_f)

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

	_, err := database.GetThreadBySlug(slug)

	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
		return
	}

	p := models.Post{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = p.UnmarshalJSON(body)

	_, err = database.GetUserInfo(p.Author)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+p.Author)
		return
	}

	existing, e := database.GetPostById(int(p.Id))
	if e == nil {
		WriteResponse(w, http.StatusConflict, existing)
		return
	}

	_, err = database.CreateThreadPost(slug, p)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, p)
}
