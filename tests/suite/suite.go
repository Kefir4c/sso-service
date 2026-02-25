package suite

import (
	"testing"

	ssov1 "github.com/Kefir4c/protos_sso/gen/go/sso"
	"github.com/Kefir4c/sso-service/internal/config"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}
