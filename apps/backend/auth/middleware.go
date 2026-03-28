package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"

// UserInfo holds the authenticated user's profile data from the JWT.
type UserInfo struct {
	ProfileID string
	Name      string
}

// Middleware extracts the Bearer token from the Authorization header,
// validates it, and stores the user info in the request context.
// It does NOT reject unauthenticated requests -- the resolver decides.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			// No "Bearer " prefix found
			next.ServeHTTP(w, r)
			return
		}

		claims, err := ValidateToken(token)
		if err != nil {
			// Invalid token -- continue without user context
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, &UserInfo{
			ProfileID: claims.ProfileID,
			Name:      claims.Name,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext retrieves the authenticated user from the context.
// Returns nil if the user is not authenticated.
func UserFromContext(ctx context.Context) *UserInfo {
	user, _ := ctx.Value(userContextKey).(*UserInfo)
	return user
}
