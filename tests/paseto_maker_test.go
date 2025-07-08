package tests

import (
	"testing"
	"time"

	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/token"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := token.NewPasetoMaker(config.RandomString(32), config.RandomString(32))
	require.NoError(t, err)

	email := config.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	payload, err := maker.CreateAccessToken(email, false, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload.Token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyAccessToken(payload.Token)
	require.NoError(t, err)
	require.NotEmpty(t, payload.Token)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(config.RandomString(32), config.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.CreateAccessToken(config.RandomEmail(), false, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload.Token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyAccessToken(payload.Token)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrExpiredToken.Error())
	require.Nil(t, payload)
}
