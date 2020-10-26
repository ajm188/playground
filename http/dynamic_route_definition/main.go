package main

import (
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	done := make(chan bool)

	go func() {
		http.ListenAndServe(":8000", nil)
		done <- true
	}()

	time.Sleep(time.Minute)
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar"))
	})

	<-done
}
