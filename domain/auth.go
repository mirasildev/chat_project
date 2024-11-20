package domain

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type AuthResponse struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	AccessToken string `json:"access_token"`
}

type VerifyTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type AuthPayload struct {
	Id        string `json:"id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	IssuedAt  string `json:"issued_at"`
	ExpiredAt string `json:"expired_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
