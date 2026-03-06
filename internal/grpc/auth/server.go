package authgrpc

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	ssov1 "github.com/Kefir4c/protos_sso/gen/go/sso"
	"github.com/Kefir4c/sso-service/internal/services/auth"
	"github.com/Kefir4c/sso-service/internal/storage"
	"github.com/Kefir4c/sso-service/internal/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, email, password string) (uid int64, err error)
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	ValidateToken(ctx context.Context, token string) (isValid bool, userID int64, email string, appID int)
	Logout(ctx context.Context, token string) (bool, error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth    Auth
	timeout time.Duration
}

const (
	errAppIDRequired  = "app_id is required"
	errUserIDRequired = "user_id is required"
	errTokenRequired  = "token is required"
)

func Register(gRPCServer *grpc.Server, auth Auth, timeout time.Duration) {
	ssov1.RegisterAuthServer(gRPCServer, &ServerAPI{auth: auth, timeout: timeout})
}

func (s *ServerAPI) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	email := strings.TrimSpace(in.GetEmail())
	password := strings.TrimSpace(in.GetPassword())

	if email == "" {
		return nil, status.Error(codes.InvalidArgument, validation.ErrEmailRequired.Error())
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, validation.ErrPasswordRequired.Error())
	}

	if err := validation.ValidateEmail(email); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidatePassword(password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	uid, err := s.auth.Register(ctx, email, password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &ssov1.RegisterResponse{UserId: uid}, nil
}

func (s *ServerAPI) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	email := strings.TrimSpace(in.GetEmail())
	password := strings.TrimSpace(in.GetPassword())

	if email == "" {
		return nil, status.Error(codes.InvalidArgument, validation.ErrEmailRequired.Error())
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, validation.ErrPasswordRequired.Error())
	}

	if err := validation.ValidateEmail(email); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidatePassword(password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, errAppIDRequired)
	}

	token, err := s.auth.Login(ctx, email, password, int(in.AppId))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login user")
	}
	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	if in.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, errUserIDRequired)
	}

	isAdmin, err := s.auth.IsAdmin(ctx, in.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func (s *ServerAPI) ValidateToken(ctx context.Context, in *ssov1.ValidateTokenRequest) (*ssov1.ValidateTokenResponse, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	if in.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, errTokenRequired)
	}

	isValid, userID, email, appID := s.auth.ValidateToken(ctx, in.GetToken())
	return &ssov1.ValidateTokenResponse{
		IsValid: isValid,
		UserId:  userID,
		Email:   email,
		AppId:   int32(appID),
	}, nil
}

func (s *ServerAPI) Logout(ctx context.Context, in *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	if in.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, errTokenRequired)
	}

	tokenPreview := in.GetToken()
	if len(tokenPreview) > 8 {
		tokenPreview = tokenPreview[:8] + "..."
	}

	slog.Info("Logout request", "Token", tokenPreview)

	success, err := s.auth.Logout(ctx, in.GetToken())
	if err != nil {
		return &ssov1.LogoutResponse{Success: false}, nil
	}
	return &ssov1.LogoutResponse{Success: success}, nil
}
