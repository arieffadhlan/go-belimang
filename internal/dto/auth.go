package dto

type (
	SignUpRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=5,max=30"`
		Password string `json:"password" validate:"required,min=5,max=30"`
	}

	SignInRequest struct {
		Username string `json:"username" validate:"required,min=5,max=30"`
		Password string `json:"password" validate:"required,min=5,max=30"`
	}

	AuthResponse struct {
		Token string `json:"token"`
	}
)
