package http

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	conv "github.com/kust1q/Zapp/main/internal/controllers/http/conv"
	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/request"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/sirupsen/logrus"
)

const (
	maxMemoryForm = 1024 * 1024 * 1024
)

// createTweet creates a new tweet for authenticated user.
//
// @Summary      Create tweet
// @Description  Create a new tweet with optional media file. Supports multipart/form-data only.
// @Tags         tweets
// @Security     Bearer
// @Accept       multipart/form-data
// @Produce      json
// @Param        content  formData  string              false  "Tweet text content"
// @Param        file     formData  file                false  "Optional media file"
// @Success      201      {object}  response.Tweet
// @Failure      400      {object}  response.Error "Invalid request body or empty tweet"
// @Failure      401      {object}  response.Error "Unauthorized"
// @Failure      500      {object}  response.Error "Internal server error"
// @Router       /protected/tweets [post]
func (h *Handler) createTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req request.Tweet
	var fileHeader *multipart.FileHeader
	var err error
	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err = c.Request.ParseMultipartForm(maxMemoryForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
			return
		}
		req.Content = c.PostForm("content")
		fileHeader, err = c.FormFile("file")
		if err != nil && err != http.ErrMissingFile {
			logrus.WithFields(logrus.Fields{
				"user_id": userID.(int),
				"error":   err,
			}).Error("create tweet failed - internal server error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	} else {
		if err := c.BindJSON(&req); err != nil {
			logrus.WithError(err).Error("failed to create tweet - invalid request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		fileHeader = nil
	}

	var file *entity.File
	if fileHeader == nil && strings.TrimSpace(req.Content) == "" {
		logrus.WithError(err).Error("failed to create tweet - impossible create empty tweet")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible create empty tweet"})
		return
	} else if fileHeader != nil {
		openedFile, err := fileHeader.Open()
		if err != nil {
			logrus.WithError(err).Error("failed to create tweet - open file error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "open file error"})
			return
		}
		defer func() {
			_ = openedFile.Close()
		}()
		file = &entity.File{
			File:   openedFile,
			Header: fileHeader,
		}
	} else {
		file = nil
	}

	tweet, err := h.tweetService.CreateTweet(c.Request.Context(), conv.FromTweetRequestToDomain(userID.(int), nil, file, &req))

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("create tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID.(int),
		"tweet_id": tweet.ID,
	}).Info("tweet created")
	c.JSON(http.StatusCreated, conv.FromDomainToTweetResponse(tweet))
}

// updateTweet updates an existing tweet of authenticated user.
//
// @Summary      Update tweet
// @Description  Update tweet content and optionally replace media. Supports multipart/form-data only.
// @Tags         tweets
// @Security     Bearer
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param        tweet_id  path      int           true   "Tweet ID"
// @Param        content   formData  string        false  "Tweet text content"
// @Param        file      formData  file          false  "Optional media file"
// @Success      200       {object}  response.Tweet
// @Failure      400       {object}  response.Error "Invalid tweet ID or request body"
// @Failure      401       {object}  response.Error "Unauthorized"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /protected/tweets/{tweet_id} [patch]
func (h *Handler) updateTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to update tweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	var req request.Tweet
	var fileHeader *multipart.FileHeader
	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err = c.Request.ParseMultipartForm(maxMemoryForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
			return
		}
		req.Content = c.PostForm("content")
		fileHeader, err = c.FormFile("file")
		if err != nil && err != http.ErrMissingFile {
			logrus.WithFields(logrus.Fields{
				"user_id": userID.(int),
				"error":   err,
			}).Error("update tweet failed - internal server error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	} else {
		if err := c.BindJSON(&req); err != nil {
			logrus.WithError(err).Error("failed to update tweet - invalid request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		fileHeader = nil
	}

	var file *entity.File
	if fileHeader == nil && strings.TrimSpace(req.Content) == "" {
		logrus.WithError(err).Error("failed to update tweet - impossible update empty tweet")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible update empty tweet"})
		return
	} else if fileHeader != nil {
		openedFile, err := fileHeader.Open()
		if err != nil {
			logrus.WithError(err).Error("failed to update tweet - open file error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "open file error"})
			return
		}
		defer func() {
			_ = openedFile.Close()
		}()
		file = &entity.File{
			File:   openedFile,
			Header: fileHeader,
		}
	} else {
		file = nil
	}

	tweet, err := h.tweetService.UpdateTweet(c.Request.Context(), conv.FromTweetUpdateRequestToDomain(userID.(int), tweetID, file, &req))
	if err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("update tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("update tweet failed - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	c.JSON(http.StatusOK, conv.FromDomainToTweetResponse(tweet))
}

// likeTweet adds like from authenticated user to tweet.
//
// @Summary      Like tweet
// @Description  Like tweet by ID for current authenticated user.
// @Tags         tweets
// @Security     Bearer
// @Produce      json
// @Param        tweet_id  path      int  true  "Tweet ID"
// @Success      200       {object}  response.Message
// @Failure      400       {object}  response.Error "Invalid tweet ID"
// @Failure      401       {object}  response.Error "Unauthorized"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /protected/tweets/{tweet_id}/like [post]
func (h *Handler) likeTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to like tweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err = h.tweetService.LikeTweet(c.Request.Context(), userID.(int), tweetID); err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("like tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("like tweet failed - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	go func() {
		if err = h.notificationService.NotifyLike(context.Background(), userID.(int), tweetID); err != nil {
			logrus.WithError(err).Error("failed to notify like")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"user_id":  userID.(int),
		"tweet_id": tweetID,
	}).Info("tweet liked")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully like tweet",
	})
}

// unlikeTweet removes like from authenticated user from tweet.
//
// @Summary      Unlike tweet
// @Description  Remove like from tweet by ID for current authenticated user.
// @Tags         tweets
// @Security     Bearer
// @Produce      json
// @Param        tweet_id  path      int  true  "Tweet ID"
// @Success      200       {object}  response.Message
// @Failure      400       {object}  response.Error "Invalid tweet ID"
// @Failure      401       {object}  response.Error "Unauthorized"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /protected/tweets/{tweet_id}/like [delete]
func (h *Handler) unlikeTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to unlike tweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err = h.tweetService.UnlikeTweet(c.Request.Context(), userID.(int), tweetID); err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("unlike tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("unlike tweet failed - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID.(int),
		"tweet_id": tweetID,
	}).Info("tweet unliked")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully unlike tweet",
	})
}

// retweet creates retweet of given tweet for authenticated user.
//
// @Summary      Retweet tweet
// @Description  Create retweet for specified tweet ID for current authenticated user.
// @Tags         tweets
// @Security     Bearer
// @Produce      json
// @Param        tweet_id  path      int  true  "Tweet ID"
// @Success      200       {object}  response.Message
// @Failure      400       {object}  response.Error "Invalid tweet ID"
// @Failure      401       {object}  response.Error "Unauthorized"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /protected/tweets/{tweet_id}/retweet [post]
func (h *Handler) retweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to retweet - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err = h.tweetService.CreateRetweet(c.Request.Context(), userID.(int), tweetID); err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("retweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("retweet failed - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	go func() {
		if err = h.notificationService.NotifyRetweet(context.Background(), userID.(int), tweetID); err != nil {
			logrus.WithError(err).Warn("failed to notify retweet")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"user_id":  userID.(int),
		"tweet_id": tweetID,
	}).Info("successfully retweet")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully retweet",
	})
}

// deleteRetweet deletes retweet created by authenticated user.
//
// @Summary      Delete retweet
// @Description  Delete existing retweet by retweet ID for current authenticated user.
// @Tags         tweets
// @Security     Bearer
// @Produce      json
// @Param        tweet_id  path      int  true  "Retweet ID"
// @Success      200       {object}  response.Message
// @Failure      400       {object}  response.Error "Invalid retweet ID"
// @Failure      401       {object}  response.Error "Unauthorized"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /protected/tweets/{tweet_id}/retweet [delete]
func (h *Handler) deleteRetweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	retweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || retweetID == 0 {
		logrus.WithError(err).Error("failed to delete retweet - invalid retweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retweet id"})
		return
	}

	if err = h.tweetService.DeleteRetweet(c.Request.Context(), userID.(int), retweetID); err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":    userID.(int),
			"retweet_id": retweetID,
			"error":      err,
		}).Error("retweet delete failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":    userID.(int),
			"retweet_id": retweetID,
			"error":      err,
		}).Error("retweet delete failed - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":    userID.(int),
		"retweet_id": retweetID,
	}).Info("successfully delete retweet")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully delete retweet",
	})
}

// replyToTweet creates reply tweet to given tweet for authenticated user.
//
// @Summary      Reply to tweet
// @Description  Create reply tweet to given tweet with optional media. Supports multipart/form-data only.
// @Tags         tweets
// @Security     Bearer
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param        tweet_id  path      int           true   "Parent tweet ID"
// @Param        content   formData  string        false  "Reply text content"
// @Param        file      formData  file          false  "Optional media file"
// @Success      201       {object}  response.Tweet
// @Failure      400       {object}  response.Error "Invalid tweet ID, body or empty reply"
// @Failure      401       {object}  response.Error "Unauthorized"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /protected/tweets/{tweet_id}/reply [post]
func (h *Handler) replyToTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	parentTweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || parentTweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	var req request.Tweet
	var fileHeader *multipart.FileHeader
	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err = c.Request.ParseMultipartForm(maxMemoryForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
			return
		}
		req.Content = c.PostForm("content")
		fileHeader, err = c.FormFile("file")
		if err != nil && err != http.ErrMissingFile {
			logrus.WithFields(logrus.Fields{
				"user_id":         userID.(int),
				"parent_tweet_id": parentTweetID,
				"error":           err,
			}).Error("reply to tweet failed - internal server error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	} else {
		if err := c.BindJSON(&req); err != nil {
			logrus.WithError(err).Error("failed to reply to tweet - invalid request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		fileHeader = nil
	}

	var file *entity.File
	if fileHeader == nil && strings.TrimSpace(req.Content) == "" {
		logrus.WithError(err).Error("failed to reply to tweet - impossible create empty tweet")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible reply to empty tweet"})
		return
	} else if fileHeader != nil {
		openedFile, err := fileHeader.Open()
		if err != nil {
			logrus.WithError(err).Error("failed to reply to tweet - open file error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "open file error"})
			return
		}
		defer func() {
			_ = openedFile.Close()
		}()
		file = &entity.File{
			File:   openedFile,
			Header: fileHeader,
		}
	} else {
		file = nil
	}

	tweet, err := h.tweetService.CreateTweet(c.Request.Context(), conv.FromTweetRequestToDomain(userID.(int), &parentTweetID, file, &req))

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":         userID.(int),
			"parent_tweet_id": parentTweetID,
			"error":           err,
		}).Error("create tweet failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	go func() {
		if err = h.notificationService.NotifyReply(context.Background(), userID.(int), parentTweetID); err != nil {
			logrus.WithError(err).Error("failed to notify reply")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"user_id":         userID.(int),
		"parent_tweet_id": parentTweetID,
	}).Info("reply to tweet created")
	c.JSON(http.StatusCreated, conv.FromDomainToTweetResponse(tweet))
}

// getReplies returns list of replies to given tweet.
//
// @Summary      Get tweet replies
// @Description  Get list of replies for specified tweet ID.
// @Tags         tweets
// @Produce      json
// @Param        tweet_id  path      int  true  "Tweet ID"
// @Success      200       {object}  response.TweetList
// @Failure      400       {object}  response.Error "Invalid tweet ID"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /public/tweets/{tweet_id}/replies [get]
func (h *Handler) getReplies(c *gin.Context) {
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
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

	replies, err := h.tweetService.GetRepliesToTweet(c.Request.Context(), tweetID, limit, offset)
	if err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get replies - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get replies - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	logrus.WithField("tweet_id", tweetID).Info("replies got")
	c.JSON(http.StatusOK, conv.FromDomainToTweetListResponse(replies))
}

// getTweetById returns tweet by ID.
//
// @Summary      Get tweet by ID
// @Description  Get single tweet by its ID.
// @Tags         tweets
// @Produce      json
// @Param        tweet_id  path      int  true  "Tweet ID"
// @Success      200       {object}  response.Tweet
// @Failure      400       {object}  response.Error "Invalid tweet ID"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /public/tweets/{tweet_id} [get]
func (h *Handler) getTweetById(c *gin.Context) {
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	tweet, err := h.tweetService.GetTweetById(c.Request.Context(), tweetID)
	if err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get tweet - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get tweet - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	logrus.WithField("tweet_id", tweetID).Info("tweet got")
	c.JSON(http.StatusOK, conv.FromDomainToTweetResponse(tweet))
}

// getTweetsAndRetweetsByUsername returns tweets and retweets for given username.
//
// @Summary      Get user tweets and retweets
// @Description  Get list of tweets and retweets for specified username.
// @Tags         tweets
// @Produce      json
// @Param        username  path      string  true  "Username"
// @Success      200       {object}  response.TweetList
// @Failure      400       {object}  response.Error "Invalid username"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /public/users/{username}/tweets [get]
func (h *Handler) getTweetsAndRetweetsByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		logrus.Error("failed to get tweet by username - invalid username")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
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

	tweets, err := h.tweetService.GetTweetsAndRetweetsByUsername(c.Request.Context(), username, limit, offset)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get tweets - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	logrus.WithField("username", username).Info("tweets got")
	c.JSON(http.StatusOK, conv.FromDomainToTweetListResponse(tweets))
}

// getLikes returns list of users who liked the tweet.
//
// @Summary      Get tweet likes
// @Description  Get list of users who liked specified tweet.
// @Tags         tweets
// @Produce      json
// @Param        tweet_id  path      int  true  "Tweet ID"
// @Success      200       {object}  response.SmallUserList
// @Failure      400       {object}  response.Error "Invalid tweet ID"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /public/tweets/{tweet_id}/likes [get]
func (h *Handler) getLikes(c *gin.Context) {
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to reply - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 50 {
		limit = 50
	}
	if limit < 1 {
		limit = 10
	}

	likes, err := h.tweetService.GetLikes(c.Request.Context(), tweetID, limit, offset)
	if err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get likes - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"error":    err,
		}).Error("failed to get likes - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	logrus.WithField("tweet_id", tweetID).Info("likes got")
	c.JSON(http.StatusOK, conv.FromDomainToSmallUserListResponse(likes))
}

// deleteTweet deletes tweet of authenticated user by ID.
//
// @Summary      Delete tweet
// @Description  Delete tweet by ID for current authenticated user.
// @Tags         tweets
// @Security     Bearer
// @Produce      json
// @Param        tweet_id  path      int  true  "Tweet ID"
// @Success      200       {object}  response.Message
// @Failure      400       {object}  response.Error "Invalid tweet ID"
// @Failure      401       {object}  response.Error "Unauthorized"
// @Failure      404       {object}  response.Error "Tweet not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /protected/tweets/{tweet_id} [delete]
func (h *Handler) deleteTweet(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok || userID.(int) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	tweetID, err := strconv.Atoi(c.Param("tweet_id"))
	if err != nil || tweetID == 0 {
		logrus.WithError(err).Error("failed to delete - invalid tweet id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet id"})
		return
	}

	if err := h.tweetService.DeleteTweet(c.Request.Context(), userID.(int), tweetID); err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("tweet delete failed - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id":  userID.(int),
			"tweet_id": tweetID,
			"error":    err,
		}).Error("tweet delete failed - tweet not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tweet not found",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userID.(int),
		"tweet_id": tweetID,
	}).Info("successfully delete tweet")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully delete tweet",
	})
}
