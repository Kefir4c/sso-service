package auth

import (
	"context"
	"time"

	ssov1 "github.com/Kefir4c/protos_sso/gen/go/sso"
)

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	timeout time.Duration
}

func (s ServerAPI) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	return nil, nil
}
