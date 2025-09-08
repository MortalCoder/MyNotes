package service

import (

	//"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"mynotes/internal/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResp struct {
	Token string `json:"token"`
}

func (s *Service) Register(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: InvalidParams})
	}

	req.Email = normalizeEmail(req.Email)
	if req.Email == "" || len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: "email or password is invalid"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("bcrypt: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}

	if err := s.usersRepo.Create(req.Email, string(hash)); err != nil {
		if dupErr(err) {
			return c.JSON(http.StatusConflict, &Response{ErrorMessage: "email already exists"})
		}
		s.logger.Errorf("usersRepo.Create: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Service) Login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: InvalidParams})
	}
	req.Email = normalizeEmail(req.Email)
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: "email or password is invalid"})
	}

	u, err := s.usersRepo.ByEmail(req.Email)
	if err != nil {
		s.logger.Errorf("usersRepo.ByEmail: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}
	if u == nil || bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(req.Password)) != nil {
		return c.JSON(http.StatusUnauthorized, &Response{ErrorMessage: "invalid credentials"})
	}

	token, err := s.signJWT(u)
	if err != nil {
		s.logger.Errorf("signJWT: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}

	return c.JSON(http.StatusOK, authResp{Token: token})
}

// ===== helpers =====

func (s *Service) signJWT(u *users.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	claims := jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

func normalizeEmail(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func dupErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique")
}
