package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/langgenius/dify-plugin-daemon/internal/service"
	"github.com/langgenius/dify-plugin-daemon/internal/types/exception"
	"github.com/langgenius/dify-plugin-daemon/pkg/validators"
)

// SubmitToolInterruptResultOutOfSession accepts {token, result} with X-Api-Key (SERVER_KEY) and forwards
// to Dify inner API; no plugin id or tenant in the body.
func SubmitToolInterruptResultOutOfSession(c *gin.Context) {
	var req service.SubmitToolInterruptOutOfSessionRequest
	if c.Request.Header.Get("Content-Type") == "application/json" {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, exception.BadRequestError(err).ToResponse())
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, exception.BadRequestError(err).ToResponse())
			return
		}
	}

	if err := validators.GlobalEntitiesValidator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, exception.BadRequestError(err).ToResponse())
		return
	}

	accepted, err := service.SubmitToolInterruptResultOutOfSession(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "plugin manager is not available" {
			c.JSON(
				http.StatusServiceUnavailable,
				exception.InternalServerError(err).ToResponse(),
			)
			return
		}
		c.JSON(
			http.StatusBadRequest,
			exception.BadRequestError(err).ToResponse(),
		)
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"accepted": accepted})
}
