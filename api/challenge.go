package api

import (
	"dionysus/models"
	"dionysus/services"
	"dionysus/controllers"

	"encoding/json"
	"net/http"
	"github.com/go-chi/chi"
	"strconv"
)

func ChallengeList(w http.ResponseWriter, r *http.Request) {
	s := services.ChallengeService{}
	challenges, err := s.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	b, err := json.Marshal(challenges)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, json.RawMessage(string(b))})
}

func ChallengeGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, "error parsing id"})
		return
	}
	s := services.ChallengeService{}
	result, err := s.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{404, "challenge does not exist"})
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

func ChallengeAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := services.ChallengeService{}
	var challenge models.Challenge
	err := json.NewDecoder(r.Body).Decode(&challenge)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	controller := controllers.ChallengeController{}
	challenge.ID, err = controller.ValidateID(challenge.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	err = s.Add(challenge)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, "Success"})
}