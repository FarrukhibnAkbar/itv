package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"itv/monorepo/api_gateway/configs"
	"itv/monorepo/api_gateway/constants"
	"itv/monorepo/api_gateway/pkg/utils"
	"net/http"
	"time"
)

// ActionLogger is a middleware for logging API requests and responses.
type ActionLogger struct {
	signingKey string
}

// NewActionLogger initializes the ActionLogger middleware.
func NewActionLogger(cfg *configs.Configuration) (*ActionLogger, error) {
	return &ActionLogger{
		signingKey: cfg.JWTSecretKey,
	}, nil
}

// Middleware logs request details, user metadata, and response status codes.
func (al *ActionLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		// Read the request body (to allow later reuse)
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		// Restore the request body
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Extract user metadata from JWT
		userMetadata, err := utils.GetUserMetadata(al.signingKey, req.Header.Get(constants.AuthorizationHeader))
		if err != nil {
			// If JWT verification fails, proceed without logging user metadata
			next.ServeHTTP(w, req)
			return
		}

		// Wrap the response writer to capture status codes
		lrw := newLoggingResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(lrw, req)

		// Extract user role and ID
		role, _ := userMetadata["role"].(string)
		id, _ := userMetadata["id"].(string)

		// Format request body for logging
		var requestBody interface{}
		if len(bodyBytes) > 0 {
			_ = json.Unmarshal(bodyBytes, &requestBody) // Attempt to parse JSON for readability
		}

		// Log request and response details
		log.Printf(`
		Role: %s | UserID: %s | Method: %s | StatusCode: %d | URL: %s | Time: %s
		Request Body: %v`,
			role, id, req.Method, lrw.statusCode, req.RequestURI, time.Since(start), requestBody)

		// Log Authorization Token on failed requests
		if lrw.statusCode >= 400 {
			log.Printf("Authorization Token: %s", req.Header.Get(constants.AuthorizationHeader))
		}
	})
}