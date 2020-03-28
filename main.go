package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/atotto/clipboard"
)

type ResponseBody struct {
	Status uint `json:"status"`
	Content string `json:"content"`
}

type RequestBody struct {
	Content string `json:"content"`
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(healthzHandler))
	mux.Handle("/halthz", http.HandlerFunc(healthzHandler))
	mux.Handle("/clipboards", http.HandlerFunc(clipboardHandler))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	body := ResponseBody{http.StatusOK, "health status is OK"}
	res, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func clipboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			clipboardGet(w, r)
		case http.MethodPost:
			clipboardPost(w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
	}
}

func clipboardGet(w http.ResponseWriter, r *http.Request) {
	body := ResponseBody{http.StatusOK, readClipboard()}
	res, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func clipboardPost(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody RequestBody
	// maximum read of 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// return a "json: unknown field"
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := ResponseBody{http.StatusOK, writeClipboard(requestBody.Content)}
	res, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func readClipboard() (string){
	t, err := clipboard.ReadAll()
	if err != nil {
		return ""
	}
	return t
}

func writeClipboard(s string) (string){
	err := clipboard.WriteAll(s)
	if err != nil {
		return ""
	}
	return readClipboard()
}
