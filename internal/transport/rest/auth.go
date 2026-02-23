package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/domain"
	"github.com/mahmud-off/auth/internal/transport/rest/response"
)

func (h *Handler) signUp(ctx *gin.Context) {

	var input domain.SignUpInput
	if err := ctx.BindJSON(&input); err != nil {
		//TODO:logging
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		//TODO:logging
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.SignUp(ctx, input); err != nil {
		//TODO:logging
		response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, response.StatusResponse{Status: "ok"})
}

func (h *Handler) signIn(ctx *gin.Context) {

	var input domain.SignInInput
	if err := ctx.BindJSON(&input); err != nil {
		//TODO: logging
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		//TODO:logging
		response.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: service to generate jwt token
	token, err := h.userService.SignIn(ctx, input)
	if err != nil {
		//TODO:logging
		response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
