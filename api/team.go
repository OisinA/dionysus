package api

import (
	"dionysus/services"

	"encoding/json"
	"net/http"
)

func TeamList(w http.ResponseWriter, r *http.Request) {
	s := services.TeamService{}
	teams, err := s.List(getQueries(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	b, err := json.Marshal(teams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, json.RawMessage(string(b))})
}

func TeamUserList(w http.ResponseWriter, r *http.Request) {
	s := services.TeamService{}
	users, err := s.ListMembers(getQueries(r))
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
