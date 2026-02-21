package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/domain"
)

type User interface {
	SignUp(ctx *gin.Context, inp domain.SignUpInput) error
}

type Handler struct {
	userService User
}

func NewHandler(users User) *Handler {
	return &Handler{userService: users}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	return router

}
