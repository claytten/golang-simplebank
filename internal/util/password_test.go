package util_test

import (
	"testing"

	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashingPassword(t *testing.T) {
	password := util.RandomString(10)
	hashedPassword, err := util.HashingPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	//reversing
	err = util.ComparePassword(hashedPassword, password)
	require.NoError(t, err)

	//testing error
	wrongPassword := util.RandomString(9)
	err = util.ComparePassword(hashedPassword, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	//testing the result is not same hash if using same raw string
	hashedPassword2, err := util.HashingPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)

	//testing hashing password if using empty string
	_, err = util.HashingPassword("")
	require.Nil(t, err)

}
