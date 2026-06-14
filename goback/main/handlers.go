package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	data := map[string]string{
		"status":      "available",
		"environment": "dev",
		"version":     "0.0.1",
	}

	json, err := json.Marshal(data)
	if err != nil {
		fmt.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getUpdateDeleteUser(w http.ResponseWriter, r *http.Request) {

}
