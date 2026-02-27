package auth_service

import (
	"context"
	"fmt"
	"net/http"

	auth_dto "github.com/Fi44er/cloud-store-api/internal/modules/auth/dto"
	auth_entity "github.com/Fi44er/cloud-store-api/internal/modules/auth/entity"
	auth_adapters "github.com/Fi44er/cloud-store-api/internal/modules/auth/infrastructure/adapters"
	auth_constant "github.com/Fi44er/cloud-store-api/internal/modules/auth/pkg/constant"
	auth_utils "github.com/Fi44er/cloud-store-api/internal/modules/auth/pkg/util"
	"github.com/Fi44er/cloud-store-api/pkg/customerr"
	getlocation "github.com/Fi44er/cloud-store-api/pkg/get_location"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/mssola/useragent"
	kratos "github.com/ory/kratos-client-go"
)

type ctxKey string

const (
	ctxKeyIP        ctxKey = "ip"
	ctxKeyUserAgent ctxKey = "ua"
)

type headerTransport struct {
	Transport http.RoundTripper
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if ip, ok := req.Context().Value(ctxKeyIP).(string); ok && ip != "" {
		req.Header.Set("X-Forwarded-For", ip)
	}
	if ua, ok := req.Context().Value(ctxKeyUserAgent).(string); ok && ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	return t.Transport.RoundTrip(req)
}

const (
	KratosPublicURL = "http://localhost:4433"
	KratosAdminURL  = "http://localhost:4434"
)

type AuthService struct {
	logger             *logger.Logger
	userUsecaseAdapter *auth_adapters.UserUsecaseAdapter

	kratosPublic *kratos.APIClient
	kratosAdmin  *kratos.APIClient
}

func NewAuthService(logger *logger.Logger) *AuthService {
	httpClient := &http.Client{
		Transport: &headerTransport{
			Transport: http.DefaultTransport,
		},
	}

	publicConfig := kratos.NewConfiguration()
	publicConfig.Servers = []kratos.ServerConfiguration{{URL: KratosPublicURL}}
	publicConfig.HTTPClient = httpClient

	adminConfig := kratos.NewConfiguration()
	adminConfig.Servers = []kratos.ServerConfiguration{{URL: KratosAdminURL}}
	adminConfig.HTTPClient = httpClient

	return &AuthService{
		logger:       logger,
		kratosPublic: kratos.NewAPIClient(publicConfig),
		kratosAdmin:  kratos.NewAPIClient(adminConfig),
	}
}

// =================== Registration ===================
func (s *AuthService) InitRegistration(ctx context.Context) (string, error) {
	flow, resp, err := s.kratosPublic.FrontendAPI.CreateNativeRegistrationFlow(ctx).Execute()

	if err != nil {
		s.logger.Errorf("failed to init registration: %v", err)
		return "", customerr.NewError(resp.StatusCode, err.Error())
	}

	return flow.Id, nil
}

func (s *AuthService) Registration(ctx context.Context, dto *auth_dto.RegistrationRequest) (*auth_dto.RegisterResponse, string, error) {
	if dto.FlowID == "" {
		s.logger.Error("flow_id is required")
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

	result, resp, err := s.kratosPublic.FrontendAPI.UpdateRegistrationFlow(ctx).
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
			Errors: errs,
		}

		s.logger.Errorf("failed to register user: %v", err)
		// TODO: вынести ошибку в константы
		return resWithErrors, "", customerr.NewError(status, "validation_failed")
	}

	user := auth_entity.User{
		ID:       result.Identity.Id,
		Email:    dto.Email,
		Username: dto.Username,
	}

	if err := s.userUsecaseAdapter.Create(ctx, &user); err != nil {
		s.logger.Errorf("Failed to create local user profile: %v", err)
		return nil, "", customerr.NewError(500, "failed to create profile")
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
		User: &auth_dto.UserDTO{
			ID:       result.Identity.Id,
			Email:    result.Identity.Traits.(map[string]any)["email"].(string),
			Username: result.Identity.Traits.(map[string]any)["username"].(string),
		},
		NeedsVerification:  needsVerification,
		VerificationFlowID: verificationFlowID,
	}

	return res, sessionToken, nil
}

// =================== Verification ===================
func (s *AuthService) InitVerification(ctx context.Context) (string, error) {
	flow, resp, err := s.kratosPublic.FrontendAPI.CreateNativeVerificationFlow(ctx).Execute()
	if err != nil {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}
		s.logger.Errorf("failed to init verification: %v", err)
		return "", customerr.NewError(status, "failed to init verification")
	}
	return flow.Id, nil
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
		request = request.Cookie(auth_constant.CratosSessionKey + "=" + cookie)
	}

	_, resp, err := request.Execute()

	if err != nil {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}

		errs := auth_utils.ExtractKratosErrors(err)
		resWithErrors := &auth_dto.VerificationResponse{
			Status: resp.Status,
			Errors: errs,
		}

		s.logger.Errorf("failed to update verification flow: %v", err)
		return resWithErrors, customerr.NewError(status, "validation_failed")
	}

	res := &auth_dto.VerificationResponse{
		Status: resp.Status,
	}

	return res, nil
}

func (s *AuthService) SendVerificationCode(ctx context.Context, flowID, email string) error {
	identities, _, err := s.kratosAdmin.IdentityAPI.ListIdentities(ctx).
		CredentialsIdentifier(email).
		Execute()

	if err == nil && len(identities) > 0 {
		identity := identities[0]
		for _, address := range identity.VerifiableAddresses {
			if address.Value == email && address.Verified {
				s.logger.Infof("email already verified")
				return customerr.NewError(400, "email_already_verified")
			}
		}
	}

	body := kratos.UpdateVerificationFlowBody{
		UpdateVerificationFlowWithCodeMethod: &kratos.UpdateVerificationFlowWithCodeMethod{
			Method: "code",
			Email:  &email,
		},
	}

	_, resp, err := s.kratosPublic.FrontendAPI.UpdateVerificationFlow(ctx).
		Flow(flowID).
		UpdateVerificationFlowBody(body).
		Execute()

	if err != nil {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}
		s.logger.Errorf("failed to update verification flow: %v", err)
		return customerr.NewError(status, err.Error())
	}
	return nil
}

// =================== Login ===================
func (s *AuthService) InitLogin(ctx context.Context) (string, error) {
	flow, resp, err := s.kratosPublic.FrontendAPI.CreateNativeLoginFlow(ctx).Execute()

	if err != nil {
		s.logger.Errorf("failed to create login flow: %v", err)
		return "", customerr.NewError(resp.StatusCode, err.Error())
	}

	return flow.Id, nil
}

func (s *AuthService) Login(ctx context.Context, dto *auth_dto.LoginRequest, userAgent, ip string) (*auth_dto.LoginResponse, string, error) {
	if dto.FlowID == "" {
		s.logger.Error("flow_id is required")
		return nil, "", customerr.NewError(400, "flow_id is required")
	}

	body := kratos.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &kratos.UpdateLoginFlowWithPasswordMethod{
			Method:     "password",
			Identifier: dto.Identifier,
			Password:   dto.Password,
		},
	}

	ctx = context.WithValue(ctx, ctxKeyIP, ip)
	ctx = context.WithValue(ctx, ctxKeyUserAgent, userAgent)

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
			Errors: errs,
		}

		s.logger.Errorf("Error: %v", err)
		return resWithErrors, "", customerr.NewError(status, "validation_failed")
	}

	isVerified := false
	if len(result.Session.Identity.VerifiableAddresses) > 0 {
		s.logger.Debugf("T: %+v", result.Session.Identity.VerifiableAddresses)
		for _, addr := range result.Session.Identity.VerifiableAddresses {
			if addr.Verified {
				isVerified = true
				break
			}
		}
	}

	if !isVerified {
		if result.SessionToken != nil {
			_, err := s.kratosPublic.FrontendAPI.PerformNativeLogout(ctx).
				PerformNativeLogoutBody(kratos.PerformNativeLogoutBody{
					SessionToken: *result.SessionToken,
				}).Execute()
			if err != nil {
				s.logger.Errorf("Error: %v", err)
				status := 500
				if resp != nil {
					status = resp.StatusCode
				}

				errs := auth_utils.ExtractKratosErrors(err)
				resWithErrors := &auth_dto.LoginResponse{
					Errors: errs,
				}

				s.logger.Errorf("Error: %v", err)
				return resWithErrors, "", customerr.NewError(status, "validation_failed")
			}
		}
		return nil, "", customerr.NewError(403, "please verify your email. link: http://localhost:8080/api/auth/verification/resend")
	}

	currentSessionID := result.Session.Id
	identityID := result.Session.Identity.Id

	sessions, _, err := s.kratosAdmin.IdentityAPI.ListIdentitySessions(ctx, identityID).Execute()
	if err == nil {
		for _, sess := range sessions {
			if sess.Id == currentSessionID || sess.Active == nil || !*sess.Active {
				continue
			}

			if len(sess.Devices) > 0 {
				oldDevice := sess.Devices[0]

				if oldDevice.IpAddress != nil && *oldDevice.IpAddress == ip &&
					oldDevice.UserAgent != nil && *oldDevice.UserAgent == userAgent {

					_, err := s.kratosAdmin.IdentityAPI.DisableSession(ctx, sess.Id).Execute()
					if err != nil {
						s.logger.Errorf("Failed to revoke duplicate device session %s: %v", sess.Id, err)
					} else {
						s.logger.Infof("Revoked old session %s for the same device (IP: %s)", sess.Id, ip)
					}
				}
			}
		}
	}

	var sessionToken string
	if result.SessionToken != nil {
		sessionToken = *result.SessionToken
	}

	res := &auth_dto.LoginResponse{
		User: &auth_dto.UserDTO{
			ID:       result.Session.Identity.Id,
			Email:    result.Session.Identity.Traits.(map[string]any)["email"].(string),
			Username: result.Session.Identity.Traits.(map[string]any)["username"].(string),
		},
		// NeedsVerification: needsVerification,
	}

	location, err := getlocation.GetIPLocation(ip)
	if err != nil {
		s.logger.Errorf("Error: %v", err)
	}

	locationRes := fmt.Sprintf("%s, %s", location.Country.Names["ru"], location.City.Names["ru"])

	patch := kratos.JsonPatch{
		Op:    "add",
		Path:  "/metadata_public/last_location",
		Value: locationRes,
	}

	_, _, err = s.kratosAdmin.IdentityAPI.PatchIdentity(ctx, identityID).
		JsonPatch([]kratos.JsonPatch{patch}).
		Execute()

	if err != nil {
		s.logger.Errorf("Failed to update identity location: %v", err)
	}

	return res, sessionToken, nil
}

// =================== Logout ===================
func (s *AuthService) Logout(ctx context.Context, cookie string) (*auth_dto.LogoutResponse, error) {
	if cookie != "" {
		resp, err := s.kratosPublic.FrontendAPI.PerformNativeLogout(ctx).
			PerformNativeLogoutBody(kratos.PerformNativeLogoutBody{SessionToken: cookie}).
			Execute()
		if err != nil {
			status := 500
			if resp != nil {
				status = resp.StatusCode
			}

			errs := auth_utils.ExtractKratosErrors(err)

			resWithErrors := &auth_dto.LogoutResponse{
				Success: false,
				Errors:  errs,
			}
			s.logger.Errorf("Error: %v", err)
			return resWithErrors, customerr.NewError(status, "logout_failed")
		}
	}

	res := &auth_dto.LogoutResponse{
		Success: true,
	}

	return res, nil
}

// =================== Sessions ===================
func (s *AuthService) GetSessions(ctx context.Context, identityID, currentSessionID string) ([]auth_dto.SessionResponse, error) {
	if identityID == "" {
		return nil, customerr.NewError(401, "no session")
	}

	sessions, _, err := s.kratosAdmin.IdentityAPI.ListIdentitySessions(context.Background(), identityID).Execute()
	if err != nil {
		return nil, err
	}

	response := s.mapSessions(sessions, currentSessionID)

	return response, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, cookie string, sessionID string) error {
	if cookie == "" {
		return customerr.NewError(401, "no session")
	}

	_, err := s.kratosAdmin.FrontendAPI.DisableMySession(ctx, sessionID).
		Cookie(auth_constant.CratosSessionKey + "=" + cookie).
		Execute()

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) mapSessions(sessions []kratos.Session, currentSessionID string) []auth_dto.SessionResponse {
	sessionsRes := make([]auth_dto.SessionResponse, 0, len(sessions))

	for _, sess := range sessions {
		if sess.Active != nil && *sess.Active {
			item := auth_dto.SessionResponse{
				ID:        sess.Id,
				IsCurrent: sess.Id == currentSessionID,
			}

			if sess.AuthenticatedAt != nil {
				item.AuthenticatedAt = *sess.AuthenticatedAt
			}

			if sess.ExpiresAt != nil {
				item.ExpiresAt = *sess.ExpiresAt
			}

			if len(sess.Devices) > 0 && sess.Devices[0].IpAddress != nil {
				item.IP = *sess.Devices[0].IpAddress

				if sess.Devices[0].UserAgent != nil {
					ua := useragent.New(*sess.Devices[0].UserAgent)
					browserName, browserVersion := ua.Browser()
					item.UserAgent = fmt.Sprintf("%s %s on %s", browserName, browserVersion, ua.OS())
				}
			}

			if sess.Identity != nil && sess.Identity.MetadataPublic != nil {
				if metadata, ok := sess.Identity.MetadataPublic.(map[string]any); ok {
					if loc, ok := metadata["last_location"].(string); ok {
						item.Location = loc
					}
				}
			}

			sessionsRes = append(sessionsRes, item)
		}
	}
	return sessionsRes
}
