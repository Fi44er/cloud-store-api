package auth_dto

import (
	kratos "github.com/ory/kratos-client-go"
)

type RegisterResponse struct {
	Success            bool              `json:"success"`
	Identity           *kratos.Identity  `json:"identity,omitempty"`
	NeedsVerification  bool              `json:"needs_verification,omitempty"`
	VerificationFlowID string            `json:"verification_flow_id,omitempty"`
	Errors             map[string]string `json:"errors,omitempty"`
}

type VerificationResponse struct {
	Success          bool                    `json:"success"`
	VerificationFlow kratos.VerificationFlow `json:"verification_flow"`
	Errors           map[string]string       `json:"errors,omitempty"`
}

type LoginResponse struct {
	Success bool              `json:"success"`
	Errors  map[string]string `json:"errors,omitempty"`
}
