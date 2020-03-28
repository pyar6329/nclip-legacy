package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/atotto/clipboard"
)

type Body struct {
	Status uint
	Content string
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/clipboards", http.HandlerFunc(clipboardHandler))
	mux.Handle("/halthz", http.HandlerFunc(healthzHandler))
	mux.Handle("/", http.HandlerFunc(healthzHandler))
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
	}
}

func readClipboard() (string){
	t, err := clipboard.ReadAll()
	if err != nil {
		return ""
	}
	return t
}
