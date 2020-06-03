package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type chatMessage struct {
	user    string
	message string
}

var srv *http.Server

var messages []chatMessage

var lock sync.Mutex

func pollHandler(w http.ResponseWriter, r *http.Request) {
	if len(messages) == 0 {
		fmt.Fprintln(w, "----- no Messages yet ---")
	} else {
		fmt.Fprint(w, "<!doctype html><html><body>")
		for _, msg := range messages {
			fmt.Fprintf(w, "<pre>%s: %s</pre>", msg.user, msg.message)
		}
		fmt.Fprint(w, "</body></html>")
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("chat.html")
	if err != nil {
		log.Fatal(err)
	} else {
		_, err := io.Copy(w, f)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	switch len(parts) {
	case 1:
		fmt.Fprintln(w, "user and message missing!")
	case 2:
		fmt.Fprintln(w, "message missing!")
	default:
		if user := parts[len(parts)-2]; user == "" {
			fmt.Fprintln(w, "empty user!")
		} else if msg := parts[len(parts)-1]; msg == "" {
			fmt.Fprintln(w, "empty message!")
		} else {
			lock.Lock()
			defer lock.Unlock()
			messages = append(messages, chatMessage{user, msg})
			fmt.Fprintln(w, "message received.")
		}
	}
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server shut down!")
	go func() {
		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func main() {
	srv = &http.Server{Addr: ":8080"}
	messages = make([]chatMessage, 0)
	http.HandleFunc("/", chatHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/poll", pollHandler)
	http.HandleFunc("/shutdown", shutdownHandler)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}
