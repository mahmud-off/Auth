package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/transport/rest/response"
)

const (
	authHeader = "Authorization"
	userCTX    = "userId"
)

func (h *Handler) authMiddleware(ctx *gin.Context) {

	token, err := getTokenFromHeader(ctx)
	if err != nil {
		//TODO: logging
		response.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	userId, err := h.userService.ParseToken(ctx, token)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set(userCTX, userId)

}

func getTokenFromHeader(ctx *gin.Context) (string, error) {

	header := ctx.GetHeader(authHeader)
	if header == "" {
		return "", errors.New("emply auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid auth header")
	}

	return headerParts[1], nil
}
