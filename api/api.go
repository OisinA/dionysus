package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"

	"dionysus/models"
	"dionysus/services"
)

type API struct{}

func (*API) Register(r chi.Router) {
	r.Get("/challenge", func(w http.ResponseWriter, r *http.Request) {
		s := services.ChallengeService{}
		challenges, err := s.List()
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		b, err := json.Marshal(challenges)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		json.NewEncoder(w).Encode(APIResponse{200, json.RawMessage(string(b))})
	})
	r.Post("/challenge", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		s := services.ChallengeService{}
		var challenge models.Challenge
		err := json.NewDecoder(r.Body).Decode(&challenge)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{400, err.Error()})
			return
		}
		err = s.Add(challenge)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{400, err.Error()})
			return
		}
		json.NewEncoder(w).Encode(APIResponse{200, "Success"})
	})
}
