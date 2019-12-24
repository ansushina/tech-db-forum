package handlers

import (
	"github.com/ansushina/tech-db-forum/pkg/database"
	"net/http"
)

func Clear(w http.ResponseWriter, r *http.Request) {
	err := database.ClearAll()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, nil)

}

func Status(w http.ResponseWriter, r *http.Request) {
	res, err := database.DatabaseStatus()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, res)
}
