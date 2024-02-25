package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/michelangelomo/ephimail/internal/redis"
)

type WebServer struct {
	Address string
	Port    int

	storage redis.Storage
}

func NewWebServer(storage redis.Storage) *WebServer {
	return &WebServer{
		storage: storage,
	}
}

func (w *WebServer) Run() error {
	m := mux.NewRouter()
	m.HandleFunc("/inbox/{email}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		emails, err := w.storage.RetrieveEmails(vars["email"])
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(emails)
	})

	fmt.Println("starting web server on", w.Address, w.Port)
	server := &http.Server{
		Handler: m,
		Addr:    fmt.Sprintf("%s:%d", w.Address, w.Port),
	}
	return server.ListenAndServe()
}
