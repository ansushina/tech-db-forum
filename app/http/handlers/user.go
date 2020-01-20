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

	us := models.Users{}
	existingUser, err := database.GetUserInfo(nickname)
	if err == nil {
		us = append(us, &existingUser)

	}
	exUser, e := database.GetUserByEmail(u.Email)
	if e == nil && existingUser != exUser {
		us = append(us, &exUser)
	}
	if err == nil || e == nil {
		WriteResponse(w, http.StatusConflict, us)
		return
	}

	var res models.User
	res, err = database.CreateUser(u)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, res)
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	nickname, _ := checkVar("nickname", r)

	res, err := database.GetUserInfo(nickname)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
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

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	_ = u.UnmarshalJSON(body)

	nickname, _ := checkVar("nickname", r)
	u.Nickname = nickname

	if u.Email != "" {
		exUser, e := database.GetUserByEmail(u.Email)
		if e == nil {
			WriteErrorResponse(w, http.StatusConflict, "This email is already registered by user: "+exUser.Nickname)
			return
		}
	}

	res, err := database.UpdateUser(u)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
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

}
