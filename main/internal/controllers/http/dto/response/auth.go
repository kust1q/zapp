package response

import "time"

type (
	SignUp struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Bio       string    `json:"bio"`
		Gen       string    `json:"gen"`
		AvatarUrl string    `json:"avatar_url"`
		CreatedAt time.Time `json:"created_at"`
	}

	Access struct {
		Access string `json:"access_token"`
	}

	Recovery struct {
		RecoveryToken string `json:"recovery_token" binding:"required,min=1,max=100"`
	}
)
