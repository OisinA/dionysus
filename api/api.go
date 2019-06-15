package api

import (
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
		})
	})
}
