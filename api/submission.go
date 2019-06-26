package api

import (
	"dionysus/services"

	"encoding/json"
	"net/http"
	"github.com/bwmarrin/lit"
	"io"
	"io/ioutil"
	"bytes"
	"github.com/go-chi/chi"
	"os"
	"strings"
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
	_ = os.Mkdir("data/submissions/" + user, os.ModePerm)
	_ = os.Mkdir("data/submissions/" + user + "/" + id, os.ModePerm)
	answer := readFile(w, r, "answer", "data/submissions/" + user + "/" + id + "/submission.txt")
	_ = readFile(w, r, "source", "data/submissions/" + user + "/" + id + "/source.txt")
	lit.Debug("Submission received.")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, "success"})
	go scoreSubmission(user, id, answer)
	return
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

func scoreSubmission(user_id string, problem_id string, answer string) {
	dat, err := ioutil.ReadFile("data/problems/" + problem_id + "/answer.txt")
	if err != nil {
		return
	}
	if strings.TrimSpace(string(dat)) == strings.TrimSpace(answer) {
		score := 10
		s := services.ScoreService{}
		_ = s.Add(user_id, problem_id, score)
	}
	return
}
