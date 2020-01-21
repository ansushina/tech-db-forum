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

	res, err := database.CreateUser(&u)
	if err == database.UserIsExist {
		WriteResponse(w, http.StatusConflict, res)
		return
	} else if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, u)
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

	res, err := database.UpdateUser(u)

	switch err {
	case nil:
		WriteResponse(w, http.StatusOK, res)
	case database.UserNotFound:
		{
			WriteErrorResponse(w, http.StatusNotFound, "Can't find user with nickname "+nickname)
			return
		}
	case database.UserConflict:
		{
			WriteErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
	default:
		{
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

}
