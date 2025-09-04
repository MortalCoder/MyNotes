package service

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (s *Service) userIDFromToken(c echo.Context) (int64, error) {
	h := c.Request().Header.Get("Authorization")
	if h == "" || !strings.HasPrefix(h, "Bearer ") {
		return 0, errors.New("no token")
	}
	raw := strings.TrimPrefix(h, "Bearer ")
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	switch v := claims["sub"].(type) {
	case float64:
		return int64(v), nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.New("bad sub")
		}
		return id, nil
	default:
		return 0, errors.New("no sub")
	}
}

func (s *Service) unauthorized(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, &Response{ErrorMessage: "unauthorized"})
}
