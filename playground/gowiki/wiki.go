package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type ChatMessage struct {
	user    string
	message string
}

var srv *http.Server

var messages []ChatMessage

func PollHandler(w http.ResponseWriter, r *http.Request) {
	if len(messages) == 0 {
		fmt.Fprintln(w, "----- no Messages yet ---")
	} else {
		for _, msg := range messages {
			fmt.Fprintf(w, "%s: %s\n", msg.user, msg.message)
		}
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {

}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	switch len(parts) {
	case 1:
		fmt.Fprintln(w, "user and message missing!")
	case 2:
		fmt.Fprintln(w, "message missing!")
	default:
		user := parts[len(parts)-2]
		msg := parts[len(parts)-1]
		messages = append(messages, ChatMessage{user, msg})
		fmt.Fprintln(w, "message received.")
	}
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server will shut down!")
	go func() {
		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func main() {
	srv = &http.Server{Addr: ":8080"}
	messages = make([]ChatMessage, 0)
	http.HandleFunc("/", ChatHandler)
	http.HandleFunc("/add/", AddHandler)
	http.HandleFunc("/poll", PollHandler)
	http.HandleFunc("/shutdown", ShutdownHandler)
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
