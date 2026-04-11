package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kust1q/Zapp/main/internal/controllers/http/conv"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

// search performs search for users or tweets by query.
//
// @Summary      Search users or tweets
// @Description  Search users or tweets by query string. Use type=user to search users; otherwise searches tweets.
// @Tags         search
// @Produce      json
// @Param        query     query     string  true   "Search query"
// @Success      200   {object}  response.SearchResult
// @Failure      400   {object}  response.Error "Missing or invalid query parameter"
// @Failure      500   {object}  response.Error "Search failed"
// @Router       /public/search [get]
func (h *Handler) search(c *gin.Context) {
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'query' is required"})
		return
	}
	users, err := h.clientSearchService.SearchUsers(c.Request.Context(), query)
	if err != nil {
		logrus.WithError(err).WithField("query", query).Error("failed to search users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	if users == nil {
		users = []entity.User{}
	}

	tweets, err := h.clientSearchService.SearchTweets(c.Request.Context(), query)
	if err != nil {
		logrus.WithError(err).WithField("query", query).Error("failed to search tweets")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	if tweets == nil {
		tweets = []entity.Tweet{}
	}

	logrus.WithField("query", query).Info("saccessfuly tweet search")
	c.JSON(http.StatusOK, gin.H{
		"tweets": conv.FromDomainToTweetListResponse(tweets),
		"users":  conv.FromDomainToUserListResponse(users),
	})
}
