package handlers

import (
	"github.com/ansushina/tech-db-forum/pkg/database"
	"net/http"
)

func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {

	slug, _ := checkVar("slug", r)

	res, err := database.GetThreadBySlug(slug)

	switch err {
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

	WriteResponse(w, http.StatusOK, res)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
