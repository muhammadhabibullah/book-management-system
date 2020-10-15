package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware check request JWT token
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		splitToken := strings.Split(bearerToken, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Forbidden: invalid token format", http.StatusForbidden)
			return
		}
		tokenString := splitToken[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Forbidden: cannot parse token", http.StatusForbidden)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Print(claims)
		} else {
			http.Error(w, "Forbidden: cannot parse token claim", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
