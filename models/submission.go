package models

type Submission struct {
	Submission_ID string    `json:"id"`
	User_ID       string    `json:"user"`
	Problem_ID    string    `json:"problem"`
	Status        int       `json:"status"`
	Score         int       `json:"score"`
	Updated       string    `json:"updated"`
}
