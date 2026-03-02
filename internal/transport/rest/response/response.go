package response

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	message string
}

type StatusResponse struct {
	Status string `json:"status"`
}

func NewErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, errorResponse{message: message})
}
