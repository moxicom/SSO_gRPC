package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moxicom/SSO_gRPC/tests/suite"
	ssov1 "github.com/moxicom/SSO_gRPC_PROTOS/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID = 1
	appSecret = "test-secret"

	passByDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, suite := suite.New(t)

	email := gofakeit.Email()
	password := randomPassword()

	responseReg, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email: email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, responseReg.GetUserId())
	
	responseLog, err := suite.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email: email,
		Pasword: password,
		AppId: appID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, responseLog.GetToken())

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(responseLog.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)
	
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, responseReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(suite.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passByDefaultLen)
}