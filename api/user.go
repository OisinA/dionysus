package api

import (
	"dionysus/services"
	"dionysus/models"

	"encoding/json"
	"net/http"
	"github.com/go-chi/chi"
)

func UserList(w http.ResponseWriter, r *http.Request) {
	s := services.UserService{}
	users, err := s.List()
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
		return
	}
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}

	user.Password, err = passwordToHash(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}

	err = s.Add(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, "success"})
}