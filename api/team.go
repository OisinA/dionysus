package api

import (
	"dionysus/services"

	"encoding/json"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/bwmarrin/lit"
)

func TeamList(w http.ResponseWriter, r *http.Request) {
	if !AdminAccess(r) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
		return
	}
	s := services.TeamService{}
	params := getQueries(r)
	teams, err := s.List(params)
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

func TeamGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s := services.TeamService{}
	result, err := s.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{404, "team does not exist"})
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

func TeamUserList(w http.ResponseWriter, r *http.Request) {
	s := services.TeamService{}
	params := getURLParams(r)
	users, err := s.ListMembers(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	if !AdminAccess(r) {
		if len(users) > 1 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
			lit.Debug("more than one team in result")
			return
		}
		user := tokenToID(r.Header.Get("Token"))
		if user == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{500, "Incorrect user."})
			lit.Debug("no user")
			return
		}
		if len(users) == 1 {
			found := false
			for _, i := range users {
				for _, j := range i {
					if j == user {
						found = true
						break
					}
				}
			}
			if !found {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
				lit.Debug("user not in team")
				return
			}
		}
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
