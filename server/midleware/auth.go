package midleware

import (
	"ChatGo/config"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type myCustomClaims struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func CreateJWT(Login string) (string, error) {

	claims := myCustomClaims{
		Login,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.Get().SigningKey))
	if err != nil {
		return "", err
	}

	return ss, nil

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenStringBearer := r.Header.Get("Authorization")
		if tokenStringBearer != "" {
			Bearer, tokenString, ok := strings.Cut(tokenStringBearer, " ")
			if ok == false || Bearer != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token, err := jwt.ParseWithClaims(tokenString, &myCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.Get().SigningKey), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "User", (token.Claims).(*myCustomClaims).User)))
			return

		} else if r.URL.Path != "/LoginUser" && r.URL.Path != "/CreateUser" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

	})
}
