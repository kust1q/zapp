package entity

import (
	"time"
)

type (
	Credential struct {
		Email    string
		Password string
	}

	Refresh struct {
		Refresh string
	}

	Access struct {
		Access string
	}

	Recovery struct {
		Recovery string
	}

	Tokens struct {
		Access  *Access
		Refresh *Refresh
	}

	UpdatePassword struct {
		UserID      int
		OldPassword string
		NewPassword string
	}

	ForgotPassword struct {
		Email string
	}

	RecoveryPassword struct {
		RecoveryToken string
		NewPassword   string
	}

	User struct {
		ID          int
		Username    string
		Gen         string
		Bio         string
		CreatedAt   time.Time
		IsSuperuser bool
		IsActive    bool
		AvatarUrl   string
		Credential  *Credential
	}

	UserProfile struct {
		User   *User
		Tweets []Tweet
	}

	Follow struct {
		FollowerID  int
		FollowingID int
		CreatedAt   time.Time
	}

	SecretQuestion struct {
		UserID         int
		SecretQuestion string
		Answer         string
	}

	SmallUser struct {
		ID        int
		Username  string
		AvatarUrl string
	}

	SecuritySettingsUpdate struct {
		UserID            int
		OldSecretAnswer   string
		NewSecretQuestion string
		NewSecretAnswer   string
	}

	UpdateBio struct {
		UserID int
		Bio    string
	}
)
