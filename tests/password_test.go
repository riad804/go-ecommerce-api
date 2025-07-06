package tests

import (
	"testing"

	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := config.RandomString(6)

	hashedPassword1, err := config.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = config.CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := config.RandomString(6)
	err = config.CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := config.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}

func TestCheckPassword(t *testing.T) {
	password := config.RandomString(10)
	hashedPassword, err := config.HashPassword(password)
	require.NoError(t, err)

	// Check with the correct password
	err = config.CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	// Check with a wrong password
	wrongPassword := config.RandomString(10)
	err = config.CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
