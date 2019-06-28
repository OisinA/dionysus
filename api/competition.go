package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/bwmarrin/lit"
)

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
		Name    string `json:"name"`
		Content string `json:"content"`
	}{name.Name, string(content)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, result})
}

func ProblemRemove(w http.ResponseWriter, r *http.Request) {
	if !AdminAccess(r) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
		return
	}
	id := chi.URLParam(r, "id")
	err := os.RemoveAll("data/problems/" + id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, "success"})
}

func ProblemAdd(w http.ResponseWriter, r *http.Request) {
	if !AdminAccess(r) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
		return
	}
	r.ParseForm()
	problem := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
		Answer  string `json:"answer"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&problem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	problem_id := uuid.New().String()
	_ = os.Mkdir("data/problems", os.ModePerm)
	_ = os.Mkdir("data/problems/"+problem_id, os.ModePerm)

	err = ioutil.WriteFile("data/problems/" + problem_id + "/answer.txt", []byte(problem.Answer), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
		return
	}

	err = ioutil.WriteFile("data/problems/" + problem_id + "/summary.md", []byte(problem.Content), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
		return
	}

	information := struct {
		Name string `json:"name"`
	}{problem.Name}

	information_marshal, err := json.Marshal(information)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
		return
	}

	err = ioutil.WriteFile("data/problems/" + problem_id + "/information.json", information_marshal, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		lit.Debug(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, "success"})
	return
}
