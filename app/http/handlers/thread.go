package handlers

import (
	"io/ioutil"
	"net/http"

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
	if err == database.ForumNotFound {
		WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	} else if err == database.UserNotFound {
		WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	} else if err == database.ThreadIsExist {
		WriteResponse(w, http.StatusConflict, res)
		return
	} else if err == nil {
		WriteResponse(w, http.StatusCreated, res)
	} else {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
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
	sort := query.Get("sort")
	if sort == "" {
		sort = "flat"
	}
	desc := query.Get("desc")
	if desc == "" {
		desc = "false"
	}

	res, err := database.GetThreadPosts(slug, limit, since, sort, desc)

	if err == database.ThreadNotFound {
		WriteErrorResponse(w, http.StatusNotFound, err.Error())
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

	_, err := database.GetThreadBySlug(slug)

	res, err := database.VoteForThread(slug, v)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.UserNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+v.Nickname)
			return
		}
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
