package entity

type (
	SmallUser struct {
		ID       int
		Username string
	}

	Tweet struct {
		ID      int
		Content string
		Author  *SmallUser
	}

	User struct {
		ID       int
		Username string
		Bio      string
	}
)
