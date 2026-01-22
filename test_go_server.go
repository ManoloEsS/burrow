package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		Handler:      mux,
	}

	mux.HandleFunc("GET /", handlerDefault)
	mux.HandleFunc("GET /health", handlerHealth)
	mux.HandleFunc("GET /monchi", handlerMonch)

	log.Fatal(s.ListenAndServe())
}

func handlerDefault(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Welcome to burrow!
You successfully started a server file and got an http response from it!
	
Keybindings are above.
You can send requests to any url, save, load and delete requests.
If you have a go http server you can write the path to the file and start it.

Try it out!
`))
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerMonch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
	w.Write([]byte("/n hello monchichi"))
}
