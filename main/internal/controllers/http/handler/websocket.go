package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// serveWs handles websocket connection requests.
//
// @Summary      Connect to WebSocket
// @Description  Upgrade HTTP connection to WebSocket for real-time notifications. Requires JWT token in query param or Authorization header.
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Param        token  query     string  false  "JWT Access Token (optional if header present)"
// @Success      101    {string}  string  "Switching Protocols"
// @Failure      401    {object}  response.Error "Unauthorized or invalid token"
// @Failure      500    {object}  response.Error "Internal server error"
// @Router       /protected/ws [get]
func (h *Handler) serveWs(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := h.webSocketService.HandleConnection(c.Writer, c.Request, userID.(int)); err != nil {
		logrus.WithError(err).Error("failed to websocket connection")
	}
}
