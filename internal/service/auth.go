package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Repository interface {
	Create(ctx *gin.Context, user domain.User) error
}

type UsersService struct {
	repo   Repository
	hasher PasswordHasher
}

func NewUsersService(repo Repository, hasher PasswordHasher) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,
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
