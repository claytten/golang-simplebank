package token_test

import (
	"testing"
	"time"

	tokens "github.com/claytten/golang-simplebank/internal/api/token"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

// add new token
func TestJWTMaker_CreateToken(t *testing.T) {
	maker, err := tokens.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	email := util.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, tokenPayload, err := maker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, tokenPayload)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

// expired token
func TestJWTMaker_ExpiredToken(t *testing.T) {
	maker, err := tokens.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomEmail(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, tokens.ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// algo none
func TestJWTMaker_InvalidAlgoNone(t *testing.T) {
	payload, err := tokens.NewPayload(util.RandomEmail(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := tokens.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, tokens.ErrInvalidToken.Error())
	require.Nil(t, payload)
}
