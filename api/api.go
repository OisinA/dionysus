package api

import (
	"github.com/go-chi/chi"
)

type API struct{}

func (*API) Register(r chi.Router) {
	// Challenges
	r.Get("/challenge", ChallengeList)
	r.Get("/challenge/{id}", ChallengeGet)
	r.Post("/challenge", ChallengeAdd)

	// Users
	r.Get("/user", UserList)
	r.Get("/user/{id}", UserGet)
	r.Post("/user", UserAdd)
}
