package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	TokenKey  contextKey = "token" // Nueva clave para rastrear el token activo
)

var (
	sessions = make(map[string]int)
	mu       sync.RWMutex
)

func CreateSession(userID int) string {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)
	mu.Lock()
	sessions[token] = userID
	mu.Unlock()
	return token
}

func DestroySession(token string) {
	mu.Lock()
	delete(sessions, token)
	mu.Unlock()
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// 1. Intentar extraer desde la Cookie (vistas HTML)
		if cookie, err := r.Cookie("session"); err == nil {
			token = cookie.Value
		}

		// 2. Si no hay cookie, intentar extraer desde el Header Authorization (Servicios Web API)
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			}
		}

		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		mu.RLock()
		userID, exists := sessions[token]
		mu.RUnlock()

		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// Almacenar tanto el ID de usuario como el token original en el contexto
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, TokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetToken recupera el token del contexto de la petición
func GetToken(ctx context.Context) string {
	if val, ok := ctx.Value(TokenKey).(string); ok {
		return val
	}
	return ""
}


func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAuthJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado"})
			return
		}
		next.ServeHTTP(w, r)
	})
}


func GetUserID(ctx context.Context) int {
	if val, ok := ctx.Value(UserIDKey).(int); ok {
		return val
	}
	return 0
}
