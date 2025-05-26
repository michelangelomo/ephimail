// server/cors.go
package server

import (
	"net/http"
	"os"
	"strings"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// NewCORSConfig creates a new CORS configuration from environment
func NewCORSConfig() *CORSConfig {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		// Default to allowing localhost in development
		appURL = "http://localhost:8080"
	}

	// Split by comma to allow multiple origins
	origins := strings.Split(appURL, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return &CORSConfig{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
}

// EnableCORS adds CORS headers to the response
func (c *CORSConfig) EnableCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// Check if origin is allowed
	allowed := false
	for _, allowedOrigin := range c.AllowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			allowed = true
			break
		}
	}

	if allowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else if len(c.AllowedOrigins) > 0 {
		// Use the first allowed origin as default
		w.Header().Set("Access-Control-Allow-Origin", c.AllowedOrigins[0])
	}

	w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.AllowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.AllowedHeaders, ", "))
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
}

// CORSMiddleware returns a middleware that handles CORS
func (c *CORSConfig) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.EnableCORS(w, r)

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
