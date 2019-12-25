package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type errorResponse struct {
	Message string `json:"message"`
}

func checkVar(varName string, req *http.Request) (string, error) {
	//log.Print(varName)
	//requestVariables := pat.Param(req, varName)
	vars := mux.Vars(req)
	return vars[varName], nil
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

	if body == nil {
		return
	}

	marshalBody, err := body.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(marshalBody)
}

func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}
