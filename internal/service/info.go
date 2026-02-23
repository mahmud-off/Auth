package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/domain"
)

type InfoRepository interface {
	GetUsers(ctx *gin.Context) ([]domain.User, error)
}

type InfoService struct {
	repo InfoRepository
}

func NewInfoService(repo InfoRepository) *InfoService {
	return &InfoService{repo: repo}
}

func (s *InfoService) GetUsers(ctx *gin.Context) ([]domain.User, error) {
	return s.repo.GetUsers(ctx)
}
