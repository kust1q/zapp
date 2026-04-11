package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authHeader = "Authorization"
	userCtx    = "userID"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Next()
		return
	}

	var token string

	header := c.GetHeader(authHeader)
	if header != "" {
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		}
	}

	if token == "" {
		token = c.Query("token")
	}

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token missing"})
		return
	}
	userID, err := h.authService.VerifyAccessToken(token)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token"})
		return
	}

	c.Set(userCtx, userID)
	c.Next()
}

func (h *Handler) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		route := c.FullPath()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			route,
		).Observe(duration)

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			route,
			strconv.Itoa(c.Writer.Status()),
		).Inc()
	}
}
