package tests

import (
	"testing"

	ssov1 "github.com/Kefir4c/protos_sso/gen/go/sso"
	"github.com/Kefir4c/sso-service/internal/lib/jwt"
	"github.com/Kefir4c/sso-service/tests/suite"
	"github.com/Kefir4c/sso-service/tests/testdata"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PasswordConfig struct {
	Length         int
	IncludeUpper   bool
	IncludeLower   bool
	IncludeNumeric bool
	IncludeSpecial bool
}

var (
	validPasswordConfig = PasswordConfig{
		Length:         10,
		IncludeUpper:   true,
		IncludeLower:   true,
		IncludeNumeric: true,
		IncludeSpecial: true,
	}

	shortPasswordConfig = PasswordConfig{
		Length:         5,
		IncludeUpper:   true,
		IncludeLower:   true,
		IncludeNumeric: true,
		IncludeSpecial: true,
	}

	longPasswordConfig = PasswordConfig{
		Length:         80,
		IncludeUpper:   true,
		IncludeLower:   true,
		IncludeNumeric: true,
		IncludeSpecial: true,
	}

	noDigitConfig = PasswordConfig{
		Length:         10,
		IncludeUpper:   true,
		IncludeLower:   true,
		IncludeNumeric: false,
		IncludeSpecial: true,
	}

	noSpecialConfig = PasswordConfig{
		Length:         10,
		IncludeUpper:   true,
		IncludeLower:   true,
		IncludeNumeric: true,
		IncludeSpecial: false,
	}

	noUpperConfig = PasswordConfig{
		Length:         10,
		IncludeUpper:   false,
		IncludeLower:   true,
		IncludeNumeric: true,
		IncludeSpecial: true,
	}
)

func generatePassword(config PasswordConfig) string {
	return gofakeit.Password(
		config.IncludeLower,
		config.IncludeUpper,
		config.IncludeNumeric,
		config.IncludeSpecial,
		false,
		config.Length,
	)
}

func TestHappyPath_RegisterLogin(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := generatePassword(validPasswordConfig)

	regResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotZero(t, regResp.GetUserId())

	logResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    testdata.AppID,
	})
	require.NoError(t, err)

	token := logResp.Token
	require.NotEmpty(t, token)

	claims, err := jwt.ValidateTokenWithSecret(token, testdata.AppSecret)
	require.NoError(t, err)

	assert.Equal(t, regResp.GetUserId(), claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, testdata.AppID, claims.AppID)
}
func TestRegister_InvalidPasswordLength(t *testing.T) {
	ctx, st := suite.New(t)

	testCases := []struct {
		name    string
		config  PasswordConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "too short (<8)",
			config:  shortPasswordConfig,
			wantErr: true,
			errMsg:  "password too short",
		},
		{
			name:    "too long (>72)",
			config:  longPasswordConfig,
			wantErr: true,
			errMsg:  "password too long",
		},
		{
			name:    "valid length (8-72)",
			config:  validPasswordConfig,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			email := gofakeit.Email()
			password := generatePassword(tc.config)

			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    email,
				Password: password,
			})
			if tc.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRegister_PasswordComplexity(t *testing.T) {
	ctx, st := suite.New(t)

	testCases := []struct {
		name    string
		config  PasswordConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no digits",
			config:  noDigitConfig,
			wantErr: true,
			errMsg:  "must contain at least one number",
		},
		{
			name:    "no special characters",
			config:  noSpecialConfig,
			wantErr: true,
			errMsg:  "must contain at least one special character",
		},
		{
			name:    "no uppercase letters",
			config:  noUpperConfig,
			wantErr: true,
			errMsg:  "must contain at least one uppercase letter",
		},
		{
			name:    "all requirements met",
			config:  validPasswordConfig,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			email := gofakeit.Email()
			password := generatePassword(tc.config)

			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    email,
				Password: password,
			})
			if tc.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRegister_Duplicate(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := generatePassword(validPasswordConfig)

	regResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotZero(t, regResp.GetUserId())

	regResp, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	require.Zero(t, regResp.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	existingEmail := gofakeit.Email()
	existingPassword := generatePassword(validPasswordConfig)

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    existingEmail,
		Password: existingPassword,
	})
	require.NoError(t, err)

	tests := []struct {
		name     string
		email    string
		password string
		appID    int32
		errMsg   string
	}{
		{
			name:     "Empty Password",
			email:    existingEmail,
			password: "",
			appID:    testdata.AppID,
			errMsg:   "password is required",
		},
		{
			name:     "Empty Email",
			email:    "",
			password: generatePassword(validPasswordConfig),
			appID:    testdata.AppID,
			errMsg:   "email is required",
		},
		{
			name:     "Both Empty",
			email:    "",
			password: "",
			appID:    testdata.AppID,
			errMsg:   "email is required",
		},
		{
			name:     "Wrong Password",
			email:    existingEmail,
			password: "Kefir_Kefr4c_Kefir",
			appID:    testdata.AppID,
			errMsg:   "invalid email or password",
		},
		{
			name:     "Non-Existent User",
			email:    "ghost@example.com",
			password: generatePassword(validPasswordConfig),
			appID:    testdata.AppID,
			errMsg:   "invalid email or password",
		},
		{
			name:     "Missing AppID",
			email:    existingEmail,
			password: existingPassword,
			appID:    0,
			errMsg:   "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.errMsg)
		})
	}
}

func TestIsAdmin(t *testing.T) {
	ctx, st := suite.New(t)

	logResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    testdata.AdminEmail,
		Password: testdata.AdminPassword,
		AppId:    testdata.AppID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, logResp.Token)

	claims, err := jwt.ValidateTokenWithSecret(logResp.GetToken(), testdata.AppSecret)
	require.NoError(t, err)

	isAdmin, err := st.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: claims.UserID,
	})
	require.NoError(t, err)
	require.True(t, isAdmin.GetIsAdmin())
}

func TestValidateToken_Success(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := generatePassword(validPasswordConfig)

	regResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotZero(t, regResp.GetUserId())

	logResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    testdata.AppID,
	})
	require.NoError(t, err)

	token := logResp.Token
	require.NotEmpty(t, token)

	validResp, err := st.AuthClient.ValidateToken(ctx, &ssov1.ValidateTokenRequest{
		Token: token,
	})
	require.NoError(t, err)
	assert.True(t, validResp.GetIsValid())
	assert.Equal(t, email, validResp.Email)
	assert.Equal(t, testdata.AppID, validResp.GetAppId())
}

func TestValidateToken_Invalid(t *testing.T) {
	ctx, st := suite.New(t)

	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty Token",
			token: "",
		},
		{
			name:  "Malformed Token",
			token: "invalid.token.string",
		},
		{
			name:  "Garbage",
			token: "gthfdrtgvcsdwertgxdsrtgfdsrtsdvbgfdgf",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validResp, err := st.AuthClient.ValidateToken(ctx, &ssov1.ValidateTokenRequest{
				Token: tc.token,
			})
			require.NoError(t, err)
			assert.False(t, validResp.GetIsValid())
		})
	}
}

func TestLogout(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := generatePassword(validPasswordConfig)

	regResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotZero(t, regResp.GetUserId())

	loginResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    testdata.AppID,
	})
	require.NoError(t, err)

	token := loginResp.GetToken()
	require.NotEmpty(t, token)

	logResp, err := st.AuthClient.Logout(ctx, &ssov1.LogoutRequest{
		Token: token,
	})
	require.NoError(t, err)
	assert.True(t, logResp.GetSuccess())

	validResp, err := st.AuthClient.ValidateToken(ctx, &ssov1.ValidateTokenRequest{
		Token: token,
	})
	require.NoError(t, err)
	assert.False(t, validResp.GetIsValid())
}
