package api

import (
	"encoding/json"
	"net/http"
	"github.com/bwmarrin/lit"
	"io/ioutil"
	"github.com/go-chi/chi"
)

var home_summary string

func CompetitionSummary(w http.ResponseWriter, r *http.Request) {
	if home_summary == "" {
		dat, err := ioutil.ReadFile("data/home_summary.md")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
			return
		}
		home_summary = string(dat)
		lit.Debug("Home summary read from file.")
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, home_summary})
	return
}

func ProblemList(w http.ResponseWriter, r *http.Request) {
	var names map[string]string = make(map[string]string, 0)
	files, err := ioutil.ReadDir("data/problems")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	for _, file := range files {
		dat, err := ioutil.ReadFile("data/problems/" + file.Name() + "/information.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
			return
		}
		name := struct {
			Name string `json:"name"`
		}{}
		err = json.Unmarshal(dat, &name)
		if err != nil {
			continue
		}
		names[file.Name()] = name.Name
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, names})
	return
}

func ProblemGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dat, err := ioutil.ReadFile("data/problems/" + id + "/information.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	name := struct {
		Name string `json:"name"`
	}{}
	err = json.Unmarshal(dat, &name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	content, err := ioutil.ReadFile("data/problems/" + id + "/summary.md")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	result := struct {
		Name string    `json:"name"`
		Content string `json:"content"`
	}{name.Name, string(content)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, result})
}