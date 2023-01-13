package util_test

import (
	"math/rand"
	"net/mail"
	"reflect"
	"testing"
	"time"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testCases := []struct {
		name     string
		min, max int64
		res      int64
	}{
		//On range
		{
			name: "MinAndMaxPositive",
			min:  1,
			max:  10000,
		},
		//Min negative
		{
			name: "MinNegativeAndMaxPositive",
			min:  -1,
			max:  10000,
		},
		//Max Negative
		{
			name: "MinPositiveAndMaxNegative",
			min:  1,
			max:  -10000,
		},
		// min max zero
		{
			name: "MinMaxZero",
			min:  0,
			max:  0,
			res:  0,
		},
		// min max range
		{
			name: "MaxRange",
			min:  -9223372036854775808,
			max:  9223372036854775807,
		},
	}

	for _, tc := range testCases {
		if tc.name == "MinAndMaxPositive" {
			t.Run(tc.name, func(t *testing.T) {
				r := util.RandomInt(tc.min, tc.max)
				require.Positive(t, r)
				require.GreaterOrEqual(t, r, tc.min)
				require.LessOrEqual(t, r, tc.max)
			})
		} else if tc.name == "MinMaxZero" {
			t.Run(tc.name, func(t *testing.T) {
				r := util.RandomInt(tc.min, tc.max)
				require.Zero(t, r)
				require.Equal(t, tc.min, r)
				require.Equal(t, tc.max, r)
			})
		} else {
			t.Run(tc.name, func(t *testing.T) {
				require.Panics(t, func() { util.RandomInt(tc.min, tc.max) }, "min and max must be positive")
			})
		}

	}
}

func TestRandomString(t *testing.T) {
	testCases := []struct {
		name string
		n    int
	}{
		//normally
		{
			name: "Random10",
			n:    10,
		},
		//normally min int
		{
			name: "RandomMinus",
			n:    -10,
		},
	}

	for _, tc := range testCases {
		if tc.name == "RandomMinus" {
			t.Run(tc.name, func(t *testing.T) {
				r := util.RandomString(tc.n)
				require.Empty(t, r)
				require.Zero(t, len(r))
			})
		} else {
			t.Run(tc.name, func(t *testing.T) {
				r := util.RandomString(tc.n)
				require.NotEmpty(t, r)
				require.Len(t, r, tc.n)
			})
		}
	}
}

func TestRandomUser(t *testing.T) {
	t.Run("RandomUser", func(t *testing.T) {
		user, password := util.RandomUser(t)
		_, err := mail.ParseAddress(user.Email)
		require.NotEmpty(t, user)
		require.NotEmpty(t, user.Username)
		require.Equal(t, len(user.Username), 50)
		require.Equal(t, len(user.FullName), 50)
		require.NoError(t, err)
		require.NotEmpty(t, password)
		require.NotZero(t, password)
		require.Equal(t, reflect.TypeOf(user).NumField(), reflect.TypeOf(db.Users{}).NumField())
	})
}

func TestRandomAccount(t *testing.T) {
	t.Run("RandomAccount", func(t *testing.T) {
		user, _ := util.RandomUser(t)
		r := util.RandomAccount(user.Username)
		require.NotEmpty(t, r)
		require.Equal(t, user.Username, r.Owner)
		require.NotEmpty(t, r.Balance)
		require.NotEmpty(t, r.Currency)
		require.NotEmpty(t, r.ID)
		require.Equal(t, reflect.TypeOf(r).NumField(), reflect.TypeOf(db.Accounts{}).NumField())
	})
}
