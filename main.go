package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
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
	s := flag.Bool("server", false, "running server")
	c := flag.Bool("copy", false, "copy from stdin")
	p := flag.Int64("port", 2230, "running port")
	flag.Usage = func() { usage() }
	flag.Parse()

	// port command: nclip --server --port 4000
	if portCheck(*p) == false {
		log.Fatal("invalid port range. Please set 1024~65535")
		os.Exit(1)
	}

	// server command: nclip --server
	if *s {
		mux := http.NewServeMux()
		mux.Handle("/", http.HandlerFunc(healthzHandler))
		mux.Handle("/healthz", http.HandlerFunc(healthzHandler))
		mux.Handle("/clipboards", http.HandlerFunc(clipboardHandler))
		ps := fmt.Sprintf(":%d", *p)
		log.Fatal(http.ListenAndServe(ps, mux))
		os.Exit(0)
	}

	// stdin command: nclip --copy
	// stdin command from pipe: echo "aaaaa" | nclip --copy
	if *c {
		if err := clipboardStdin(p); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// stdout command: nclip
	st, err := clipboardStdout(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(st)
	os.Exit(0)
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

func clipboardStdout(port *int64) (string, error) {
	s, err := clipboardGetClient(port)
	if err != nil {
		return readClipboard(), nil
	}
	return s, nil
}

func clipboardStdin(port *int64) error {
	// from stdin
	if terminal.IsTerminal(syscall.Stdin) {
		input := bufio.NewScanner(os.Stdin)
		var text []string
		for input.Scan() {
			text = append(text, input.Text())
		}
		s := strings.Join(text, "\n")
		if err := clipboardPostClient(port, s); err != nil {
			writeClipboard(s)
		}
		return nil
	}
	// from pipe
	return clipboardStdinFromPipe(port)
}

func clipboardStdinFromPipe(port *int64) error {
	s, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	if err := clipboardPostClient(port, string(s)); err != nil {
		writeClipboard(string(s))
	}
	return nil
}

func clipboardGetClient(port *int64) (string, error) {
	host := fmt.Sprintf("127.0.0.1:%d", *port)
	url := &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "clipboards",
	}
	// timeout is 1 second
	client := &http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println("Clipboard Send was failed. Please check port option")
		return "", err
	}
	defer res.Body.Close()
	var responseBody ResponseBody
	if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return "", err
	}
	return responseBody.Content, nil
}

func clipboardPostClient(port *int64, s string) error {
	host := fmt.Sprintf("127.0.0.1:%d", *port)
	url := &url.URL{
		Scheme: "http",
		Host:   host,
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
		log.Println("Clipboard Send was failed. Please check port option")
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

func portCheck(p int64) bool {
	if p >= 1024 && p <= 65535 {
		return true
	}
	return false
}
