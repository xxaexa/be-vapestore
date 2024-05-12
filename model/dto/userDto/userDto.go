package userDto

type (
	LoginUserRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	CreateUserRequest struct {
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"fullname" binding:"required"`
		Password string `json:"password" binding:"required,min=8,max=20"`
	}

	UpdateUserRequest struct {
		ID       string `json:"id" binding:"required"`
		FullName string `json:"fullname" binding:"required"`
		Password string `json:"password" binding:"required,min=8,max=20"`
	}
)
