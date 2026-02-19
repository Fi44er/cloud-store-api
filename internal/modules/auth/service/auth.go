package auth_service

import (
	"context"

	auth_dto "github.com/Fi44er/cloud-store-api/internal/modules/auth/dto"
	auth_utils "github.com/Fi44er/cloud-store-api/internal/modules/auth/pkg/util"
	"github.com/Fi44er/cloud-store-api/pkg/customerr"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	kratos "github.com/ory/kratos-client-go"
)

const (
	KratosPublicURL = "http://localhost:4433"
	KratosAdminURL  = "http://localhost:4434"
)

type AuthService struct {
	logger *logger.Logger

	kratosPublic *kratos.APIClient
	kratosAdmin  *kratos.APIClient
}

func NewAuthService(logger *logger.Logger) *AuthService {
	publicConfig := kratos.NewConfiguration()
	publicConfig.Servers = []kratos.ServerConfiguration{
		{URL: KratosPublicURL},
	}

	adminConfig := kratos.NewConfiguration()
	adminConfig.Servers = []kratos.ServerConfiguration{
		{URL: KratosAdminURL},
	}
	return &AuthService{
		logger:       logger,
		kratosPublic: kratos.NewAPIClient(publicConfig),
		kratosAdmin:  kratos.NewAPIClient(adminConfig),
	}
}

func (s *AuthService) InitRegistration(ctx context.Context) (string, error) {
	flow, resp, err := s.kratosPublic.FrontendAPI.CreateNativeRegistrationFlow(ctx).Execute()

	if err != nil {
		return "", customerr.NewError(resp.StatusCode, err.Error())
	}

	return flow.Id, nil
}

func (s *AuthService) Registration(ctx context.Context, dto *auth_dto.RegistrationRequest) (*auth_dto.RegisterResponse, string, error) {
	if dto.FlowID == "" {
		return nil, "", customerr.NewError(400, "flow_id is required")
	}

	body := kratos.UpdateRegistrationFlowBody{
		UpdateRegistrationFlowWithPasswordMethod: &kratos.UpdateRegistrationFlowWithPasswordMethod{
			Method:   "password",
			Password: dto.Password,
			Traits: map[string]any{
				"email":    dto.Email,
				"username": dto.Username,
			},
		},
	}

	result, resp, err := s.kratosPublic.FrontendAPI.UpdateRegistrationFlow(context.Background()).
		Flow(dto.FlowID).
		UpdateRegistrationFlowBody(body).
		Execute()

	if err != nil {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}

		errs := auth_utils.ExtractKratosErrors(err)

		resWithErrors := &auth_dto.RegisterResponse{
			Success: false,
			Errors:  errs,
		}

		return resWithErrors, "", customerr.NewError(status, "validation_failed")
	}

	var sessionToken string
	if result.SessionToken != nil {
		sessionToken = *result.SessionToken
	}

	var verificationFlowID string
	needsVerification := false
	if result.ContinueWith != nil {
		for _, item := range result.ContinueWith {
			if item.ContinueWithVerificationUi != nil {
				needsVerification = true
				verificationFlowID = item.ContinueWithVerificationUi.Flow.Id
			}
		}
	}

	res := &auth_dto.RegisterResponse{
		Success:            true,
		Identity:           &result.Identity,
		NeedsVerification:  needsVerification,
		VerificationFlowID: verificationFlowID,
	}

	return res, sessionToken, nil
}

func (s *AuthService) Verification(ctx context.Context, dto *auth_dto.VerificationRequest, cookie string) (*auth_dto.VerificationResponse, error) {
	body := kratos.UpdateVerificationFlowBody{
		UpdateVerificationFlowWithCodeMethod: &kratos.UpdateVerificationFlowWithCodeMethod{
			Method: "code",
			Code:   &dto.Code,
		},
	}

	request := s.kratosPublic.FrontendAPI.UpdateVerificationFlow(ctx).
		Flow(dto.FlowID).
		UpdateVerificationFlowBody(body)

	if cookie != "" {
		request = request.Cookie("ory_kratos_session=" + cookie)
	}

	result, resp, err := request.Execute()

	if err != nil {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}

		errs := auth_utils.ExtractKratosErrors(err)
		resWithErrors := &auth_dto.VerificationResponse{
			Success: false,
			Errors:  errs,
		}

		return resWithErrors, customerr.NewError(status, "validation_failed")
	}

	res := &auth_dto.VerificationResponse{
		Success:          true,
		VerificationFlow: *result,
	}

	return res, nil
}

func (s *AuthService) InitLogin(ctx context.Context) (string, error) {
	flow, resp, err := s.kratosPublic.FrontendAPI.CreateNativeLoginFlow(ctx).Execute()

	if err != nil {
		return "", customerr.NewError(resp.StatusCode, err.Error())
	}

	return flow.Id, nil
}

func (s *AuthService) Login(ctx context.Context, dto *auth_dto.LoginRequest) (*auth_dto.LoginResponse, string, error) {
	if dto.FlowID == "" {
		return nil, "", customerr.NewError(400, "flow_id is required")
	}

	body := kratos.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &kratos.UpdateLoginFlowWithPasswordMethod{
			Method:     "password",
			Identifier: dto.Identifier,
			Password:   dto.Password,
		},
	}

	result, resp, err := s.kratosPublic.FrontendAPI.UpdateLoginFlow(ctx).
		Flow(dto.FlowID).
		UpdateLoginFlowBody(body).
		Execute()

	if err != nil {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}

		errs := auth_utils.ExtractKratosErrors(err)

		resWithErrors := &auth_dto.LoginResponse{
			Success: false,
			Errors:  errs,
		}

		return resWithErrors, "", customerr.NewError(status, "validation_failed")
	}

	var sessionToken string
	if result.SessionToken != nil {
		sessionToken = *result.SessionToken
	}

	res := &auth_dto.LoginResponse{
		Success: true,
	}

	return res, sessionToken, nil
}
