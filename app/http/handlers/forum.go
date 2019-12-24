package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/ansushina/tech-db-forum/pkg/database"
	"goji.io/pat"
	"io/ioutil"
	"net/http"
	"strconv"
)

type errorResponse struct {
	Message string `json:"message"`
}

func checkVar(varName string, req *http.Request) (string, error) {
	requestVariables := pat.Param(req, varName)

	return requestVariables, nil
}

func WriteErrorResponse(w http.ResponseWriter, errCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)
	marshalBody, err := json.Marshal(errorResponse{Message: errMsg})
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(marshalBody)
}

func WriteResponse(w http.ResponseWriter, code int, body interface{ MarshalJSON() ([]byte, error) }) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	marshalBody, err := body.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(marshalBody)
}

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	f := models.Forum{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = f.UnmarshalJSON(body)

	_, err := database.GetUserInfo(f.User)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+f.User)
		return
	}

	existingForum, err := database.GetForumBySlug(f.Slug)
	if err == nil {
		WriteResponse(w, http.StatusConflict, existingForum)
		return
	}

	_, err = database.CreateForum(f)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, f)
	w.WriteHeader(http.StatusOK)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	slug, _ := checkVar("slug", r)
	f, err := database.GetForumBySlug(slug)
	switch err {
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
	/*jsonBlob, e := result.MarshalJSON()
	if e != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}*/

	WriteResponse(w, http.StatusOK, f)
	w.WriteHeader(http.StatusOK)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	slug, _ := checkVar("slug", r)
	query := r.URL.Query()
	limit := query.Get("limit")
	since := query.Get("since")
	desc, _ := strconv.ParseBool(query.Get("desc"))

	_, e := database.GetForumBySlug(slug)
	switch e {
	case database.ForumNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find forum with slug"+slug)
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, e.Error())
			return
		}
	}

	res, err := database.GetForumThreads(slug, limit, since, desc)

	switch err {
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

	/*jsonBlob, e := res.MarshalJSON()
	if e != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusOK, jsonBlob)*/

	w.WriteHeader(http.StatusOK)
}

func ForumGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
