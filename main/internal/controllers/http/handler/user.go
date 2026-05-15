package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	conv "github.com/kust1q/Zapp/main/internal/controllers/http/conv"
	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/request"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/sirupsen/logrus"
)

// getMe returns profile of authenticated user.
//
// @Summary      Get current user profile
// @Description  Get profile information of the currently authenticated user.
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Success      200  {object}  response.UserProfile
// @Failure      401  {object}  response.Error "Unauthorized"
// @Failure      404  {object}  response.Error "User not found"
// @Failure      500  {object}  response.Error "Internal server error"
// @Router       /protected/users/me [get]
func (h *Handler) getMe(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
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

	userProfile, err := h.userService.GetMe(c.Request.Context(), userID.(int), limit, offset)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to get user profile - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to get user profile - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	logrus.WithField("user_id", userID.(int)).Info("successfully get profile")
	c.JSON(http.StatusOK, conv.FromDomainToUserProfileResponse(userProfile))
}

// updateMe updates bio of authenticated user.
//
// @Summary      Update current user bio
// @Description  Update profile bio of the currently authenticated user.
// @Tags         users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        request  body      request.UpdateBio  true  "Update bio data"
// @Success      200      {object}  response.Message
// @Failure      400      {object}  response.Error "Invalid request body"
// @Failure      401      {object}  response.Error "Unauthorized"
// @Failure      404      {object}  response.Error "User not found"
// @Failure      500      {object}  response.Error "Internal server error"
// @Router       /protected/users/me [patch]
func (h *Handler) updateMe(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	var req request.UpdateBio
	if err := c.BindJSON(&req); err != nil {
		logrus.WithError(err).Error("failed to update user bio - invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.userService.Update(c.Request.Context(), conv.FromUpdateBioRequestToDomain(userID.(int), &req)); err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to update bio - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to update bio - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	logrus.WithField("user_id", userID.(int)).Info("successfully update bio")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully upate bio",
	})
}

// deleteMe deletes authenticated user account.
//
// @Summary      Delete current user
// @Description  Delete account of the currently authenticated user.
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Success      200  {object}  response.Message
// @Failure      401  {object}  response.Error "Unauthorized"
// @Failure      404  {object}  response.Error "User not found"
// @Failure      500  {object}  response.Error "Internal server error"
// @Router       /protected/users/me [delete]
func (h *Handler) deleteMe(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID.(int)); err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to delete user - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"user_id": userID.(int),
			"error":   err,
		}).Error("failed to delete user - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	logrus.WithField("user_id", userID.(int)).Info("successfully delete user")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully delete user",
	})
}

// followers returns list of followers for given username.
//
// @Summary      Get followers
// @Description  Get list of followers for specified username.
// @Tags         users
// @Produce      json
// @Param        username  path      string  true  "Username"
// @Success      200       {object}  response.SmallUserList
// @Failure      400       {object}  response.Error "Invalid username"
// @Failure      404       {object}  response.Error "User not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /public/users/{username}/followers [get]
func (h *Handler) followers(c *gin.Context) {
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

	followers, err := h.userService.GetFollowers(c.Request.Context(), username, limit, offset)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get followers - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get followers - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	logrus.WithField("username", username).Info("successfully get followers")
	c.JSON(http.StatusOK, conv.FromDomainToSmallUserListResponse(followers))
}

// following returns list of followings for given username.
//
// @Summary      Get followings
// @Description  Get list of users that specified username is following.
// @Tags         users
// @Produce      json
// @Param        username  path      string  true  "Username"
// @Success      200       {object}  response.SmallUserList
// @Failure      400       {object}  response.Error "Invalid username"
// @Failure      404       {object}  response.Error "User not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /public/users/{username}/following [get]
func (h *Handler) following(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		logrus.Error("failed to get tweet by username - invalid username")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
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

	followings, err := h.userService.GetFollowings(c.Request.Context(), username, limit, offset)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get followings - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get followings - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	logrus.WithField("username", username).Info("successfully get followers")
	c.JSON(http.StatusOK, conv.FromDomainToSmallUserListResponse(followings))
}

// followUser subscribes authenticated user to another user.
//
// @Summary      Follow user
// @Description  Follow user by ID as current authenticated user.
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Param        user_id  path      int  true  "User ID to follow"
// @Success      200      {object}  response.Follow
// @Failure      400      {object}  response.Error "Invalid user ID"
// @Failure      401      {object}  response.Error "Unauthorized"
// @Failure      404      {object}  response.Error "User not found"
// @Failure      500      {object}  response.Error "Internal server error"
// @Router       /protected/users/{user_id}/follow [post]
func (h *Handler) followUser(c *gin.Context) {
	followerID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	followingID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || followingID == 0 {
		logrus.WithError(err).Error("failed follow - invalid user id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	follow, err := h.userService.FollowToUser(c.Request.Context(), followerID.(int), followingID)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"follower_id":  followerID,
			"following_id": followingID,
			"error":        err,
		}).Error("failed to follow - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"follower_id":  followerID,
			"following_id": followingID,
			"error":        err,
		}).Error("failed to follow - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	go func() {
		if err = h.notificationService.NotifyFollow(context.Background(), followerID.(int), followingID); err != nil {
			logrus.WithError(err).Error("failed to notify follow")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"follower_id":  followerID,
		"following_id": followingID,
	}).Info("successfully follow")
	c.JSON(http.StatusOK, conv.FromDomainToFollow(follow))
}

// unfollowUser unsubscribes authenticated user from another user.
//
// @Summary      Unfollow user
// @Description  Unfollow user by ID as current authenticated user.
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Param        user_id  path      int  true  "User ID to unfollow"
// @Success      200      {object}  response.Message
// @Failure      400      {object}  response.Error "Invalid user ID"
// @Failure      401      {object}  response.Error "Unauthorized"
// @Failure      404      {object}  response.Error "User not found"
// @Failure      500      {object}  response.Error "Internal server error"
// @Router       /protected/users/{user_id}/follow [delete]
func (h *Handler) unfollowUser(c *gin.Context) {
	followerID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	followingID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || followingID == 0 {
		logrus.WithError(err).Error("failed unfollow - invalid user id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userService.UnfollowUser(c.Request.Context(), followerID.(int), followingID); err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"follower_id":  followerID,
			"following_id": followingID,
			"error":        err,
		}).Error("failed to unfollow - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"follower_id":  followerID,
			"following_id": followingID,
			"error":        err,
		}).Error("failed to unfollow - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	logrus.WithFields(logrus.Fields{
		"follower_id":  followerID,
		"following_id": followingID,
	}).Info("successfully unfollow")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully unfollow",
	})
}

// getUserProfile returns public profile by username.
//
// @Summary      Get user profile
// @Description  Get public profile information for specified username.
// @Tags         users
// @Produce      json
// @Param        username  path      string  true  "Username"
// @Success      200       {object}  response.UserProfile
// @Failure      400       {object}  response.Error "Invalid username"
// @Failure      404       {object}  response.Error "User not found"
// @Failure      500       {object}  response.Error "Internal server error"
// @Router       /public/users/{username}/profile [get]
func (h *Handler) getUserProfile(c *gin.Context) {
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

	profile, err := h.userService.GetUserProfile(c.Request.Context(), username, limit, offset)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get user profile - internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if errors.Is(err, errs.ErrUserNotFound) {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to get user profile - user not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	logrus.WithField("username", username).Info("successfully get profile")
	c.JSON(http.StatusOK, conv.FromDomainToUserProfileResponse(profile))
}
