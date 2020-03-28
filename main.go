package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ResponseBody struct {
	Status  uint   `json:"status"`
	Content string `json:"content"`
}

type RequestBody struct {
	Content string `json:"content"`
}

func usage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	s := flag.Bool("server", false, "running server default port: 8080")
	flag.Usage = func() { usage() }
	flag.Parse()

	switch *s {
	case true:
		mux := http.NewServeMux()
		mux.Handle("/", http.HandlerFunc(healthzHandler))
		mux.Handle("/halthz", http.HandlerFunc(healthzHandler))
		mux.Handle("/clipboards", http.HandlerFunc(clipboardHandler))
		log.Fatal(http.ListenAndServe(":8080", mux))
	default:
		// fmt.Print(clipboardGetClient())
		if err := clipboardPostClient("🤔"); err != nil {
			fmt.Print(err)
		}
	}
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

func clipboardGetClient() string {
	url := &url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   "clipboards",
	}
	// timeout is 1 second
	client := &http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return ""
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	var responseBody ResponseBody
	if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return ""
	}
	return responseBody.Content
}

func clipboardPostClient(s string) error {
	url := &url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   "clipboards",
	}
	// timeout is 1 second
	client := &http.Client{Timeout: time.Duration(1) * time.Second}
	requestBody := new(bytes.Buffer)
	if err := json.NewEncoder(requestBody).Encode(&RequestBody{Content: s}); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), requestBody)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var responseBody ResponseBody
	if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return err
	}
	return nil
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
	// maximum read of 1MB
	var maximumBodySize int64 = 1048576
	var requestBody RequestBody

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maximumBodySize)

	// return a "json: unknown field"
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&requestBody); err != nil {
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

func readClipboard() string {
	t, err := clipboard.ReadAll()
	if err != nil {
		return ""
	}
	return t
}

func writeClipboard(s string) string {
	if err := clipboard.WriteAll(s); err != nil {
		return ""
	}
	return readClipboard()
}
