package request

type (
	Tweet struct {
		Content string `json:"content" binding:"required,min=1,max=280"`
	}
)
