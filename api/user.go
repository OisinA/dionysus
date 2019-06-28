package api

import (
	"dionysus/services"
	"dionysus/models"

	"encoding/json"
	"net/http"
	"github.com/go-chi/chi"
	"fmt"
	"github.com/bwmarrin/lit"
)

func UserList(w http.ResponseWriter, r *http.Request) {
	if !AdminAccess(r) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
		return
	}
	s := services.UserService{}
	users, err := s.List(getQueries(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	b, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, json.RawMessage(string(b))})
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !AdminAccess(r) {
		if id != tokenToID(r.Header.Get("Token")) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
			return
		}
	}
	s := services.UserService{}
	result, err := s.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{404, "user does not exist"})
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, json.RawMessage(string(b))})
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := services.UserService{}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Error(err.Error())
		return
	}

	role := 0
	token := r.Header.Get("Token")
	if token != "" {
		if IsAdmin(token) {
			role = user.Role
		}
	}
	lit.Debug(fmt.Sprint(role))
	user.Role = role

	user.Password, err = passwordToHash(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Error(err.Error())
		return
	}

	err = s.Add(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, "success"})
	lit.Debug("New user added: " + user.Username)
	lit.Debug(fmt.Sprint(user))
}