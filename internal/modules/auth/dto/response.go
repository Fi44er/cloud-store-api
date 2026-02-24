package auth_dto

import (
	"time"
)

type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type RegisterResponse struct {
	User               *UserDTO          `json:"user,omitempty"`
	NeedsVerification  bool              `json:"needs_verification"`
	VerificationFlowID string            `json:"verification_flow_id,omitempty"`
	Errors             map[string]string `json:"errors,omitempty"`
}

type VerificationResponse struct {
	Status string            `json:"status,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

type LoginResponse struct {
	User              *UserDTO          `json:"user,omitempty"`
	NeedsVerification bool              `json:"needs_verification"`
	Errors            map[string]string `json:"errors,omitempty"`
}

type LogoutResponse struct {
	Success bool              `json:"success"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type SessionResponse struct {
	ID              string    `json:"id"`
	IP              string    `json:"ip"`
	UserAgent       string    `json:"user_agent"`
	Location        string    `json:"location"`
	IsCurrent       bool      `json:"is_current"`
	AuthenticatedAt time.Time `json:"authenticated_at"`
	ExpiresAt       time.Time `json:"expires_at"`
}
