// server/web_socket_integration.go
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/michelangelomo/ephimail/internal/redis"
)

// WebServerWithWebSocket extends WebServer with WebSocket capabilities
type WebServerWithWebSocket struct {
	*WebServer
	wsHub *WebSocketHub
}

// NewWebServerWithWebSocket creates a new web server with WebSocket support
func NewWebServerWithWebSocket(storage redis.Storage, domains []string) *WebServerWithWebSocket {
	return &WebServerWithWebSocket{
		WebServer: NewWebServer(storage, domains),
		wsHub:     NewWebSocketHub(),
	}
}

// Run starts the web server with WebSocket support
func (w *WebServerWithWebSocket) Run() error {
	// Start WebSocket hub
	go w.wsHub.Run()

	// Create router
	m := mux.NewRouter()

	// Apply CORS middleware to all routes
	m.Use(w.corsConfig.CORSMiddleware)

	// Add WebSocket endpoint
	m.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		ServeWs(w.wsHub, rw, r)
	})

	// Config endpoint for frontend runtime configuration
	m.HandleFunc("/api/config", func(rw http.ResponseWriter, r *http.Request) {
		config := map[string]interface{}{
			"backendUrl": os.Getenv("BACKEND_URL"),
		}

		// If BACKEND_URL is not set, use the request's origin
		if config["backendUrl"] == nil || config["backendUrl"] == "" {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			config["backendUrl"] = fmt.Sprintf("%s://%s", scheme, r.Host)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(config)
	})

	// Add existing endpoints
	m.HandleFunc("/domains", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(w.domains)
	})

	m.HandleFunc("/inbox/{email}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		emails, err := w.storage.RetrieveEmails(vars["email"])
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(emails)
	})

	// Register reservation handlers
	w.RegisterReservationHandlers(m)

	// Serve static files
	staticPath := "./frontend/dist"
	staticFileDirectory := http.Dir(staticPath)
	staticFileHandler := http.StripPrefix("/dist/", http.FileServer(staticFileDirectory))
	m.PathPrefix("/dist/").Handler(staticFileHandler)

	// SPA catch-all handler
	m.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticPath+"/index.html")
	})

	server := &http.Server{
		Handler: m,
		Addr:    fmt.Sprintf("%s:%d", w.Address, w.Port),
	}

	return server.ListenAndServe()
}

// GetWebSocketHub returns the WebSocket hub
func (w *WebServerWithWebSocket) GetWebSocketHub() *WebSocketHub {
	return w.wsHub
}

// NotifyNewEmail notifies WebSocket clients about a new email
func (w *WebServerWithWebSocket) NotifyNewEmail(to, messageID string) {
	w.wsHub.NotifyNewEmail(to, messageID)
}

// MailServerWithWebSocket extends MailServer with WebSocket notifications
type MailServerWithWebSocket struct {
	*MailServer
	webServer *WebServerWithWebSocket
}

// NewMailServerWithWebSocket creates a new mail server with WebSocket notifications
func NewMailServerWithWebSocket(storage redis.Storage, webServer *WebServerWithWebSocket) *MailServerWithWebSocket {
	return &MailServerWithWebSocket{
		MailServer: NewMailServer(storage),
		webServer:  webServer,
	}
}

// Run starts the mail server with WebSocket integration
func (m *MailServerWithWebSocket) Run() {
	storage, ok := m.storage.(*redis.RedisStorage)
	if !ok {
		// Fallback to standard mail server if type assertion fails
		m.MailServer.Run()
		return
	}

	// Run with encryption and WebSocket support
	m.MailServer.RunWithEncryption(storage, m.webServer.GetWebSocketHub())
}
