package token_test

import (
	"testing"
	"time"

	tokens "github.com/claytten/golang-simplebank/internal/api/token"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

// add new paseto
func TestNewPasetoMaker(t *testing.T) {
	maker, err := tokens.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	email := util.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	pasetoToken, pasetoTokenPayload, err := maker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, pasetoToken)
	require.NotEmpty(t, pasetoTokenPayload)

	payload, err := maker.VerifyToken(pasetoToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

// expired token paseto
func TestExpiredPasetoMaker(t *testing.T) {
	maker, err := tokens.NewPasetoMaker(util.RandomString(32))
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
