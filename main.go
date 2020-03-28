package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/atotto/clipboard"
)

type Body struct {
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
	body := Body{http.StatusOK, "health status is OK"}
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
			body := Body{http.StatusOK, readClipboard()}
			res, err := json.Marshal(body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		case http.MethodPost:
			if r.Header.Get("Content-Type") != "application/json" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var requestBody RequestBody
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			body := Body{http.StatusOK, writeClipboard(requestBody.Content)}
			res, err := json.Marshal(body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		default:
			w.WriteHeader(http.StatusBadRequest)
	}
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
