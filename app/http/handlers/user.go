package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/ansushina/tech-db-forum/app/models"
	"github.com/ansushina/tech-db-forum/pkg/database"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = u.UnmarshalJSON(body)

	nickname, _ := checkVar("nickname", r)
	u.Nickname = nickname

	existingUser, err := database.GetUserInfo(nickname)
	if err == nil {
		WriteResponse(w, http.StatusConflict, existingUser)
		return
	}

	var res models.User
	res, err = database.CreateUser(u)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, res)
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	nickname, _ := checkVar("nickname", r)

	res, err := database.GetUserInfo(nickname)

	switch err {
	case database.UserNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+nickname)
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

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = u.UnmarshalJSON(body)

	nickname, _ := checkVar("nickname", r)
	u.Nickname = nickname

	res, err := database.UpdateUser(u)

	switch err {
	case database.UserNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+nickname)
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
