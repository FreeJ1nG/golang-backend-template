package dto

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh"`
}

type SignUpRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type SignUpResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh"`
}

type GetCurrentUserResponse struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh"`
}
