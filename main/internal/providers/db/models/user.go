package models

import (
	"time"
)

type (
	User struct {
		ID          int       `db:"id"`
		Username    string    `db:"username"`
		Email       string    `db:"email"`
		Password    string    `db:"password"`
		Bio         string    `db:"bio"`
		Gen         string    `db:"gen"`
		CreatedAt   time.Time `db:"created_at"`
		IsActive    bool      `db:"is_active"`
		IsSuperuser bool      `db:"is_superuser"`
	}

	Follow struct {
		FollowerID  int       `db:"follower_id"`
		FollowingID int       `db:"following_id"`
		CreatedAt   time.Time `db:"created_at"`
	}
)
