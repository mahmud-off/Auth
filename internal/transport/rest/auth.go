package rest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/domain"
	"github.com/mahmud-off/auth/internal/transport/rest/response"
	"github.com/mahmud-off/auth/pkg/logger"
)

func (h *Handler) signUp(ctx *gin.Context) {

	var input domain.SignUpInput
	if err := ctx.BindJSON(&input); err != nil {
		logger.Errorf("JSON Unmarshalling problem %s", err.Error())
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		logger.Errorf("JSON data validation problem: %s", err.Error())
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.SignUp(ctx, input); err != nil {
		logger.Errorf("Signing-Up problem: %s", err.Error())
		response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, response.StatusResponse{Status: "ok"})
}

func (h *Handler) signIn(ctx *gin.Context) {

	var input domain.SignInInput
	if err := ctx.BindJSON(&input); err != nil {
		logger.Errorf("JSON Unmarshalling problem %s", err.Error())
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		logger.Errorf("JSON data validation problem: %s", err.Error())
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.userService.SignIn(ctx, input)
	if err != nil {
		logger.Errorf("Signing-In problem: %s", err.Error())
		if err == sql.ErrNoRows {
			logger.Debug("Trying to sign-in as unexisted user")
			response.NewErrorResponse(ctx, http.StatusUnauthorized, "trying to sign-in as unexisted user")
		} else {
			logger.Debug("Internal error")
			response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	//ctx.SetCookie("refresh-token", refreshToken, 0, "", "", true, true)
	ctx.Writer.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))

	ctx.Writer.Header().Add("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}

func (h *Handler) refresh(ctx *gin.Context) {

	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		logger.Errorf("Invalid refresh token: %s", err.Error())
		response.NewErrorResponse(ctx, http.StatusBadRequest, "invalid refresh token")
		return
	}

	accessToken, refreshToken, err := h.userService.RefreshTokens(ctx, cookie)

	if err != nil {
		logger.Errorf("Refresh tokens problem: %s", err.Error())
		response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	//ctx.SetCookie("refresh-token", refreshToken, 0, "", "", true, true)
	ctx.Writer.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))

	ctx.Writer.Header().Add("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}
