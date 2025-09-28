package middleware

import (
	"belimang/internal/utils"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthContext struct {
	ID      string
	IsAdmin bool
}

func Protected(requireAdmin bool) func(http.Handler) http.Handler {
	secret := []byte(os.Getenv("JWT_SECRET"))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.SendErrorResponse(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid bearer token")
				return
			}

			tknStr := parts[1]
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			if err != nil || !token.Valid {
				utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid or expired JWT")
				return
			}

			id, _ := claims["id"].(string)
			isAdmin, _ := claims["is_admin"].(bool)

			if requireAdmin && !isAdmin {
				utils.SendErrorResponse(w, http.StatusForbidden, "forbidden")
				return
			}
			if !requireAdmin && isAdmin {
				utils.SendErrorResponse(w, http.StatusForbidden, "forbidden")
				return
			}

			authCtx := AuthContext{ID: id, IsAdmin: isAdmin}
			ctx := context.WithValue(r.Context(), AuthContext{}, authCtx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthContext(ctx context.Context) (AuthContext, bool) {
	val := ctx.Value(AuthContext{})
	if val == nil {
		return AuthContext{}, false
	}

	ac, ok := val.(AuthContext)
	return ac, ok
}
