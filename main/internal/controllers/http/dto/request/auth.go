package request

type (
	SignUp struct {
		Username string `json:"username" binding:"required,alphanum,max=64"`
		Email    string `json:"email" binding:"required,email,max=64"`
		Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
		Gen      string `json:"gen" binding:"required,oneof=male female"`
		Bio      string `json:"bio" binding:"max=140"`
	}

	SignIn struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
	}

	Refresh struct {
		Refresh string `json:"refresh_token" binding:"required"`
	}

	UpdatePassword struct {
		OldPassword string `json:"old_password" binding:"required,alphanum,min=8,max=64"`
		NewPassword string `json:"new_password" binding:"required,alphanum,min=8,max=64"`
	}

	ForgotPassword struct {
		Email string `json:"email" binding:"required,email,max=64"`
	}

	RecoveryPassword struct {
		RecoveryToken string `json:"recovery_token" binding:"required,min=8,max=100"`
		NewPassword   string `json:"new_password" binding:"required,alphanum,min=8,max=64"`
	}
)
