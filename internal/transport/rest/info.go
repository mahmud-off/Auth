package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahmud-off/auth/internal/transport/rest/response"
)

func (h *Handler) GetUsers(ctx *gin.Context) {

	users, err := h.infoService.GetUsers(ctx)
	if err != nil {
		//TODO: logging
		response.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}
