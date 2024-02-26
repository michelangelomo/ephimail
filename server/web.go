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
	domains []string
}

func NewWebServer(storage redis.Storage, domains []string) *WebServer {
	return &WebServer{
		storage: storage,
		domains: domains,
	}
}

func (w *WebServer) SetAllowedDomains(domains []string) {
	w.domains = domains
}

func (w *WebServer) Run() error {
	m := mux.NewRouter()
	m.HandleFunc("/domains", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(w.domains)
	})

	m.HandleFunc("/inbox/{email}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		emails, err := w.storage.RetrieveEmails(vars["email"])
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(emails)
	})

	// Serve static files
	staticPath := "./frontend/dist"
	staticFileDirectory := http.Dir(staticPath)
	staticFileHandler := http.StripPrefix("/dist/", http.FileServer(staticFileDirectory))
	m.PathPrefix("/dist/").Handler(staticFileHandler)

	// SPA catch-all handler
	m.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticPath+"/index.html")
	})

	fmt.Println("starting web server on", w.Address, w.Port)
	server := &http.Server{
		Handler: m,
		Addr:    fmt.Sprintf("%s:%d", w.Address, w.Port),
	}
	return server.ListenAndServe()
}
