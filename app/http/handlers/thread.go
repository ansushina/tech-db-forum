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

	_, e := database.GetForumBySlug(slug)
	if e == database.ForumNotFound {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find forum with slug"+slug)
		return
	} else if e != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, e.Error())
		return
	}

	_, err := database.GetUserInfo(t.Author)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+t.Author)
		return
	}

	if t.Slug != "" {
		existing, existErr := database.GetThreadBySlug(t.Slug)
		//fmt.Println(existing)
		if existErr == nil {
			WriteResponse(w, http.StatusConflict, existing)
			return
		}
	}

	t.Forum = slug

	res, err := database.CreateForumThread(t)

	if err == nil {
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

	th, e := database.GetThreadBySlug(slug)
	if e == database.ThreadNotFound {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
		return
	} else if e != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, e.Error())
		return

	}

	res, err := database.GetThreadPosts(strconv.Itoa(th.Id), limit, since, sort, desc)
	//log.Print(res)

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

	_, err := database.GetThreadBySlug(slug)

	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
		return
	}

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

	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find thread with slug "+slug)
		return
	}

	_, err = database.GetUserInfo(v.Nickname)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+v.Nickname)
		return
	}

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
