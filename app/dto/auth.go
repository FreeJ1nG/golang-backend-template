package dto

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Token string `json:"token"`
}

type GetCurrentUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"fullName"`
}
