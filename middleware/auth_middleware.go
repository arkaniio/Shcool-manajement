package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware/logger"
	"github.com/ArkaniLoveCoding/Shcool-manajement/utils"
)

// LoggerResponseWriter wraps http.ResponseWriter to capture status code
type LoggerResponseWriter struct {
	statusCode int
	http.ResponseWriter
}

// WriteHeader captures the status code
func (lwr *LoggerResponseWriter) WriteHeader(code int) {
	lwr.statusCode = code
	lwr.ResponseWriter.WriteHeader(code)
}

// Write captures the bytes written (for detecting client disconnection)
func (lwr *LoggerResponseWriter) Write(b []byte) (int, error) {
	n, err := lwr.ResponseWriter.Write(b)
	if err != nil {
		// Client disconnected (socket hang up) - log with available info
		logger.Log.Warn("Client disconnected during response write",
			zap.Error(err),
		)
	}
	return n, err
}

// LoggerResponse middleware logs all HTTP requests with detailed information
// including client disconnection detection
func LoggerResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := &LoggerResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logger.Log.Info("HTTP Request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("ip", r.RemoteAddr),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", duration),
		)
	})
}

// getRequestID retrieves or generates a request ID from the request
func GetRequestIDInternal(r *http.Request) string {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = r.Header.Get("X-Requestid")
	}
	return requestID
}

// RequestIDMiddleware adds a request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(r *http.Request) string {
	if reqID := r.Context().Value("request_id"); reqID != nil {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	return ""
}

// TokenIdMiddleware middleware validates JWT token and extracts user info
func TokenIdMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get request ID for logging
		requestID := GetRequestID(r)

		// Get authorization header
		header := r.Header.Get("Authorization")
		if header == "" {
			logger.Log.Warn("Missing authorization header",
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
				zap.String("path", r.URL.Path),
			)
			utils.ResponseError(w, http.StatusBadRequest, "The header of the token is nil!", false)
			return
		}

		// Extract token from Bearer prefix
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			logger.Log.Warn("Empty token after Bearer prefix",
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "The token is nil!", false)
			return
		}

		// Validate the token
		token_validate, err := utils.ValidateToken(token)
		if err != nil {
			logger.Log.Warn("Token validation failed",
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
				zap.Error(err),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the token!", err.Error())
			return
		}
		if token_validate == nil {
			logger.Log.Warn("Token validation returned nil",
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the data validate!", false)
			return
		}

		// Parse user ID from token
		user_id, err := uuid.Parse(token_validate.Id)
		if err != nil {
			logger.Log.Warn("Failed to parse UUID from token",
				zap.String("request_id", requestID),
				zap.Error(err),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to convert into an uuid!", err.Error())
			return
		}

		// Save user ID to context
		user_id_ctx := context.WithValue(r.Context(), "user_id", user_id)
		if user_id_ctx == nil {
			logger.Log.Error("Failed to create context with user_id",
				zap.String("request_id", requestID),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the context value in request client!", false)
			return
		}
		r = r.WithContext(user_id_ctx)

		// Save user role to context
		role_user_ctx := context.WithValue(r.Context(), "role_user", token_validate.Role)
		if role_user_ctx == nil {
			logger.Log.Error("Failed to create context with role_user",
				zap.String("request_id", requestID),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the context value in request client", false)
			return
		}
		r = r.WithContext(role_user_ctx)

		// Log successful authentication
		logger.Log.Debug("User authenticated successfully",
			zap.String("request_id", requestID),
			zap.String("user_id", user_id.String()),
			zap.String("role", token_validate.Role),
		)

		// Continue to next handler
		next.ServeHTTP(w, r)

	})
}

// GetIdMiddleware retrieves user ID from context
func GetIdMiddleware(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {

	// Get token from context
	user_id := r.Context().Value("user_id")
	if user_id == "" && user_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token jwt!", false)
		return uuid.Nil, nil
	}

	// Convert to UUID
	uuid_user, ok := user_id.(uuid.UUID)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert from string into an uuid!", ok)
		return uuid.Nil, nil
	}
	if uuid_user == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to load the uuid user!", false)
		return uuid.Nil, nil
	}

	//return final result
	return uuid_user, nil

}

// GetRoleMiddleware retrieves user role from context
func GetRoleMiddleware(w http.ResponseWriter, r *http.Request) (string, error) {

	// Get from context
	role_context := r.Context().Value("role_user")
	if role_context == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role user from context!", false)
		return "", nil
	}

	// Parse to string
	role_user, ok := role_context.(string)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert data from type any into a string type data!", false)
		return "", nil
	}
	if role_user == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the role user from context!", false)
		return "", nil
	}

	//return final result
	return role_user, nil

}

