package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/Sigaeasu/go-mwe/config"
	"github.com/Sigaeasu/go-mwe/utils"
	"github.com/Sigaeasu/go-mwe/utils/response"
	"github.com/golang-jwt/jwt/v4"
)

type key int

const (
	Customer key = iota
)

func AuthMiddlewareService() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtConfig := config.Config.JWTCfg
			authorizationHeader := r.Header.Get("Authorization")
			url_to_skip_auth_check := []string{"/api/v1/init",}
			skip_check := utils.Contains(r.URL.Path, url_to_skip_auth_check)
			if skip_check {
				next.ServeHTTP(w, r)
				return
			}
			if !strings.Contains(authorizationHeader, "Token") {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.ResponseAPI{
					Error_: &response.ApiError{
						Error: "Invalid Token",
					},
				})
				return
			}
			tokenString := strings.Replace(authorizationHeader, "Token ", "", -1)

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Signing Method Invalid")
				} else if method != JWT_SIGNING_METHOD {
					return nil, fmt.Errorf("Signing Method Invalid")
				}

				return []byte(jwtConfig.SignKey), nil
			})
			if err != nil {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.ResponseAPI{
					Error_: &response.ApiError{
						Error: err.Error(),
					},
				})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.ResponseAPI{
					Error_: &response.ApiError{
						Error: err.Error(),
					},
				})
				return
			}

			ctxt := context.WithValue(r.Context(), Customer, claims)
			r = r.WithContext(ctxt)
			next.ServeHTTP(w, r)
		})
	}
}
