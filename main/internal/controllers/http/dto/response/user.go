package response

import "time"

type (
	SmallUser struct {
		Username  string `json:"username"`
		AvatarUrl string `json:"avatar_url"`
	}

	User struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Bio       string    `json:"bio"`
		Gen       string    `json:"gen"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		AvatarUrl string    `json:"avatar_url"`
	}

	UserProfile struct {
		User   *User   `json:"user"`
		Tweets []Tweet `json:"tweets"`
	}

	Follow struct {
		FollowerID  int       `json:"follower_id"`
		FollowingID int       `json:"following_id"`
		CreatedAt   time.Time `json:"created_at"`
	}

	// For docs
	SmallUserList struct {
		Users []SmallUser `json:"users"`
	}
)
