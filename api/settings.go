package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bwmarrin/lit"
)

type Settings struct {
	Name     string `json:"name"`
	Homepage string `json:"homepage"`
}

var settings Settings

func LoadSettings() error {
	if _, err := os.Stat("data/settings.json"); os.IsNotExist(err) {
		input, err := ioutil.ReadFile("templates/settings.json")
		if err != nil {
			lit.Error("Templates missing!")
			return err
		}
		err = ioutil.WriteFile("data/settings.json", input, 0644)
		if err != nil {
			lit.Error("Could not create settings file.")
			return err
		}
	}
	settingsf, err := ioutil.ReadFile("data/settings.json")
	if err != nil {
		return err
	}
	settings_file := struct {
		Name     string `json:"Name"`
		Homepage string `json:"Homepage"`
	}{}
	json.Unmarshal([]byte(settingsf), &settings_file)
	settings = Settings{settings_file.Name, settings_file.Homepage}
	return nil
}

func writeSettings() {
	s, err := json.Marshal(settings)
	if err != nil {
		lit.Error("Could not marshal settings json.")
		return
	}
	err = ioutil.WriteFile("data/settings.json", s, 0644)
	if err != nil {
		lit.Error("Could not create settings file.")
		return
	}
}

func CompetitionSummary(w http.ResponseWriter, r *http.Request) {
	if settings.Homepage == "" {
		err := LoadSettings()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
			return
		}
		lit.Debug("Settings read from file.")
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, settings.Homepage})
	return
}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	if !AdminAccess(r) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "no permission"})
		return
	}
	r.ParseForm()
	s := Settings{}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	settings = s
	writeSettings()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, "success"})
	return
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	if settings.Name == "" {
		err := LoadSettings()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
			return
		}
		lit.Debug("Settings read from file.")
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{200, settings})
	return
}
