package api

import (
	"github.com/go-chi/chi"
)

type API struct{}

func (*API) Register(r chi.Router) {
	// Authentication
	r.Post("/login", Login)

	r.Group(func(r chi.Router) {
		r.Use(authenticationHandler)

		// Challenges
		r.Get("/challenge", ChallengeList)
		r.Get("/challenge/{id}", ChallengeGet)
		r.Post("/challenge", ChallengeAdd)

		// Users
		r.Get("/user", UserList)
		r.Get("/user/{id}", UserGet)
		r.Post("/user", UserAdd)
	})
}
