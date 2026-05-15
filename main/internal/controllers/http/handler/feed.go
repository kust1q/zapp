package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	conv "github.com/kust1q/Zapp/main/internal/controllers/http/conv"
	"github.com/sirupsen/logrus"
)

// getFeed returns feed for authenticated user.
//
// @Summary      Get user feed
// @Description  Get tweets feed for current authenticated user (subscriptions, own tweets, etc.).
// @Tags         feed
// @Security     Bearer
// @Produce      json
// @Success      200  {object}  response.TweetList
// @Failure      401  {object}  response.Error "Unauthorized"
// @Failure      500  {object}  response.Error "Internal server error"
// @Router       /protected/feed [get]
func (h *Handler) getFeed(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 30 {
		limit = 30
	}
	if limit < 1 {
		limit = 10
	}

	feed, err := h.feedService.GetUserFeedByUserId(c.Request.Context(), userID.(int), limit, offset)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to get feed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.WithField("user_id", userID.(int)).Info("successfully get feed")
	c.JSON(http.StatusOK, conv.FromDomainToTweetListResponse(feed))
}

// getFeed returns feed for everybody.
//
// @Summary      Get user feed
// @Description  Get tweets feed for any user.
// @Tags         feed
// @Produce      json
// @Success      200  {object}  response.TweetList
// @Failure      500  {object}  response.Error "Internal server error"
// @Router       /public [get]
func (h *Handler) getDefaultFeed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 30 {
		limit = 30
	}
	if limit < 1 {
		limit = 10
	}

	feed, err := h.feedService.GetDeafultFeed(c.Request.Context(), limit, offset)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get default feed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	logrus.Info("successfully get feed")
	c.JSON(http.StatusOK, conv.FromDomainToTweetListResponse(feed))
}
