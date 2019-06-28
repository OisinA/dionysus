package api

import (
	"dionysus/services"

	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"fmt"
	"time"

	"github.com/bwmarrin/lit"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var discord_url string = ""

const (
	status_submitted = iota
	status_scoring
	status_scored
	status_failed
)

func SubmissionAdd(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := tokenToID(r.Header.Get("Token"))
	if user == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, "Incorrect user."})
		lit.Debug("no user")
		return
	}
	uuid := uuid.New()
	_ = os.Mkdir("data/submissions/"+user, os.ModePerm)
	_ = os.Mkdir("data/submissions/"+user+"/"+id, os.ModePerm)
	_ = os.Mkdir("data/submissions/"+user+"/"+id + "/" + uuid.String(), os.ModePerm)
	answer := readFile(w, r, "answer", "data/submissions/"+user+"/"+id+"/" + uuid.String() + "/submission.txt")
	_ = readFile(w, r, "source", "data/submissions/"+user+"/"+id+"/" + uuid.String() + "/source.txt")
	lit.Debug("Submission received.")

	s := services.SubmissionService{}
	s.Add(uuid.String(), user, id, status_submitted)
	//sendWebhook("Submission by user " + user + " on problem " + id + ".")
	go scoreSubmission(uuid.String(), user, id, answer)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, uuid.String()})
	return
}

func SubmissionList(w http.ResponseWriter, r *http.Request) {
	user := tokenToID(r.Header.Get("Token"))
	params := getQueries(r)
	if !AdminAccess(r) {
		if id, ok := params.Queries["User_ID"]; ok {
			if id != user {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
				return
			}
		} else {
			params.Queries["User_ID"] = user
		}
	}
	s := services.SubmissionService{}
	submissions, err := s.List(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	for i, _ := range submissions {
		if submissions[i].Status == 2 {
			s := services.ScoreService{}
			score, err := s.GetSubmissionScore(submissions[i].Submission_ID)
			if err == nil {
				submissions[i].Score = score
			}
		}
	}
	b, err := json.Marshal(submissions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, json.RawMessage(string(b))})
}

func SubmissionLastUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s := services.SubmissionService{}
	params := services.SearchParams{make(map[string]interface{}, 0)}
	params.Queries["Submission_ID"] = id
	lit.Debug(fmt.Sprint(params))
	submissions, err := s.List(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	lit.Debug(fmt.Sprint(submissions))
	b, err := json.Marshal(submissions[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, json.RawMessage(string(b))})
}

func sendWebhook(message string) {
	j := struct {
		Content  string `json:"content"`
		Username string `json:"username"`
	}{message, "dionysus"}
	b, err := json.Marshal(j)
	if err != nil {
		lit.Error(err.Error())
		return
	}
	lit.Debug(string(b))
	_, err = http.Post("https://discordapp.com/api/webhooks/593788139844141066/C4G0HgmAMwNwV_7pfaYVXkg0XEB0kS0DQQwchpuTL_6aqTD3rlqlQo6Fv5jexOJbb167", "application/json", bytes.NewBuffer(b))
	if err != nil {
		lit.Error(err.Error())
		return
	}
}

func GetScores(w http.ResponseWriter, r *http.Request) {
	user := tokenToID(r.Header.Get("Token"))
	if user == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, "Incorrect user."})
		lit.Debug("no user")
		return
	}
	s := services.ScoreService{}
	scores, err := s.Get(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, scores})
}

func GetSubmissionScore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s := services.ScoreService{}
	scores, err := s.GetSubmissionScore(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, scores})
}

func readFile(w http.ResponseWriter, r *http.Request, filename string, filepath string) string {
	var buf bytes.Buffer
	file, _, err := r.FormFile(filename)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
		return ""
	}
	io.Copy(&buf, file)
	err = ioutil.WriteFile(filepath, buf.Bytes(), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
		return ""
	}
	return buf.String()
}

func scoreSubmission(submission_id string, user_id string, problem_id string, answer string) {
	time.Sleep(30000 * time.Millisecond)
	dat, err := ioutil.ReadFile("data/problems/" + problem_id + "/answer.txt")
	s := services.SubmissionService{}
	s.Add(submission_id, user_id, problem_id, status_scoring)
	if err != nil {
		return
		s.Add(submission_id, user_id, problem_id, status_failed)
	}
	score := 0
	if strings.TrimSpace(string(dat)) == strings.TrimSpace(answer) {
		score = 10
		//sendWebhook("User scored " + fmt.Sprint(score) + " points for problem " + problem_id + ".")
	}
	s1 := services.ScoreService{}
	_ = s1.Add(submission_id, user_id, problem_id, score)
	s.Add(submission_id, user_id, problem_id, status_scored)
	return
}
