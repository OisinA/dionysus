package api

import (
	"dionysus/services"

	"net/http"
	"encoding/json"
	"github.com/go-chi/chi"
)

type API struct{}

func (*API) Register(r chi.Router) {
	// Authentication
	r.Group(func(r chi.Router) {
		r.Post("/login", Login)

		r.Get("/token_to_id", TokenToID)
		r.Post("/user", UserAdd)

		r.Group(func(r chi.Router) {
			r.Use(AuthenticationHandler)

			// Challenges
			r.Get("/challenge", ChallengeList)
			r.Get("/challenge/{id}", ChallengeGet)
			r.Post("/challenge", ChallengeAdd)

			// Users
			r.Get("/user", UserList)
			r.Get("/user/{id}", UserGet)

			// Teams
			r.Get("/team", TeamList)
			r.Get("/team_members", TeamUserList)
		})
	})
}

func getQueries(r *http.Request) services.SearchParams {
	r.ParseForm()
	params := services.SearchParams{make(map[string]interface{}, 0)}
	json.NewDecoder(r.Body).Decode(&params.Queries)
	return params
}
