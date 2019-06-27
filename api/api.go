package api

import (
	"dionysus/services"

	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type API struct{}

func (*API) Register(r chi.Router) {
	// Authentication
	r.Group(func(r chi.Router) {
		r.Post("/login", Login)

		r.Get("/settings", GetSettings)

		r.Get("/token_to_id", TokenToID)
		r.Post("/user", UserAdd)

		r.Get("/home_summary", CompetitionSummary)

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
			r.Get("/team/{id}", TeamGet)
			r.Get("/team_members", TeamUserList)

			// Problems
			r.Get("/problem", ProblemList)
			r.Get("/problem/{id}", ProblemGet)

			// Submissions
			r.Post("/submission/{id}", SubmissionAdd)
			r.Get("/scores", GetScores)
			r.Get("/submission", SubmissionList)
			r.Get("/submission/{id}", SubmissionLastUpdate)
			r.Get("/submission/{id}/score", GetSubmissionScore)

			// Settings
			r.Post("/settings", UpdateSettings)
		})
	})
}

func getQueries(r *http.Request) services.SearchParams {
	r.ParseForm()
	params := services.SearchParams{make(map[string]interface{}, 0)}
	json.NewDecoder(r.Body).Decode(&params.Queries)
	return params
}

func getURLParams(r *http.Request) services.SearchParams {
	r.ParseForm()
	m := make(map[string]interface{}, 0)
	for i := range r.URL.Query() {
		m[i] = r.URL.Query()[i][0]
	}
	params := services.SearchParams{m}
	return params
}
