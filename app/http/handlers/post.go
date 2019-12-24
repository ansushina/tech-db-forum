package handlers

import (
	"net/http"
	"strconv"
	"github.com/ansushina/tech-db-forum/pkg/database"
)

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id_str, _ := checkVar("id", r)
	id, _ := strconv.Atoi(id_str)

	res, err := database.GetPostById(id)

	switch err {
	case database.PostNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find post with id " + id_str)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	WriteResponse(w, http.StatusOK, res)

	w.WriteHeader(http.StatusOK)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
