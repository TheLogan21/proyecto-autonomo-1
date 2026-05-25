package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

type contextKey string

const UserIDKey contextKey = "userID"

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
		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		mu.RLock()
		userID, exists := sessions[cookie.Value]
		mu.RUnlock()

		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func GetUserID(ctx context.Context) int {
	if val, ok := ctx.Value(UserIDKey).(int); ok {
		return val
	}
	return 0
}
