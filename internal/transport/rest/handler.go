package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/domain"
)

type User interface {
	SignUp(ctx *gin.Context, inp domain.SignUpInput) error
	SignIn(ctx *gin.Context, inp domain.SignInInput) (string, string, error)
	LogOut(ctx *gin.Context, accessToken string, refreshToken string) error
	TokenInBlackList(ctx *gin.Context, token string) bool
	ParseToken(cxt *gin.Context, token string) (int, error)

	RefreshTokens(ctx *gin.Context, refreshToken string) (string, string, error)
}

type Info interface {
	GetUsers(ctx *gin.Context) ([]domain.User, error)
}

type Handler struct {
	userService User
	infoService Info
}

func NewHandler(users User, info Info) *Handler {
	return &Handler{
		userService: users,
		infoService: info,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
		auth.GET("/log-out", h.logOut).Use(h.authMiddleware)
		auth.GET("/refresh", h.refresh)
	}

	info := router.Group("/info")
	info.Use(h.authMiddleware)
	{
		info.GET("/", h.GetUsers)
	}

	return router

}
