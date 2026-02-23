package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mahmud-off/auth/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Repository interface {
	Create(ctx *gin.Context, user domain.User) error
	GetByCredentials(cxt *gin.Context, email string, password string) (domain.User, error)
}

type UsersService struct {
	repo   Repository
	hasher PasswordHasher

	hmacSecret []byte
}

func NewUsersService(repo Repository, hasher PasswordHasher, secret []byte) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,

		hmacSecret: secret,
	}
}

func (s *UsersService) SignUp(ctx *gin.Context, inp domain.SignUpInput) error {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIn(cxt *gin.Context, inp domain.SignInInput) (string, error) {

	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByCredentials(cxt, inp.Email, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:  strconv.Itoa(user.ID),
		IssuedAt: time.Now().Unix(),
		//TODO: link from config --> .yml
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	return token.SignedString(s.hmacSecret)

}

func (s *UsersService) ParseToken(ctx *gin.Context, token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return s.hmacSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		//TODO: logging
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return id, nil
}
