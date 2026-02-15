package auth

import (
	"context"
	"time"

	ssov1 "github.com/Kefir4c/protos_sso/gen/go/sso"
	"github.com/Kefir4c/sso-service/internal/validation"
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
	timeout time.Duration
	auth    Auth
}

func ValidateAuth(in *ssov1.AuthClient)

func (s ServerAPI) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	if err := validation.ValidateEmail(in.Email); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidatePassword(in.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// uid, err := s.auth.Register(ctx, in.GetEmail(), in.GetPassword())

	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
