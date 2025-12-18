package auth

type Create struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type PhoneRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type GetOTP struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
	Token string `json:"token"`
}
type CacheOTP struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type GetInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Token     string `json:"token"`
}
