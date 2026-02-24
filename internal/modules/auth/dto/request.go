package auth_dto

type RegistrationRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	FlowID   string `json:"flow_id"`
}

type VerificationRequest struct {
	FlowID string `json:"flow_id"`
	Code   string `json:"code"`
}

type VerificationSendCodeRequest struct {
	FlowID string `json:"flow_id"`
	Email  string `json:"email"`
}

type LoginRequest struct {
	FlowID     string `json:"flow_id"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
