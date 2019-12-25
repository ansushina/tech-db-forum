package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/ansushina/tech-db-forum/pkg/database"
)

func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	t := models.Thread{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = t.UnmarshalJSON(body)
	slug, _ := checkVar("slug", r)

	t.Forum = slug

	res, err := database.CreateForumThread(t)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.ThreadNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {

	slug, _ := checkVar("slug", r)

	res, err := database.GetThreadBySlug(slug)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.ThreadNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	slug, _ := checkVar("slug", r)
	query := r.URL.Query()
	limit := query.Get("limit")
	since := query.Get("since")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	_, e := database.GetThreadBySlug(slug)
	if e == database.ForumNotFound {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
		return
	} else if e != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, e.Error())
		return

	}

	res, err := database.GetThreadPosts(slug, limit, since, desc)

	if err == database.ThreadNotFound {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
		return
	} else if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return

	}
	WriteResponse(w, http.StatusOK, res)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	t := models.Thread{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = t.UnmarshalJSON(body)
	slug, _ := checkVar("slug", r)

	t.Forum = slug

	res, err := database.UpdateThreadBySlugorId(slug, t)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.ThreadNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	v := models.Vote{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = v.UnmarshalJSON(body)
	slug, _ := checkVar("slug", r)

	res, err := database.VoteForThread(slug, v)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.ThreadNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
