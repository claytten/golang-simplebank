package util

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/stretchr/testify/require"
)

const Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	if min < 0 || max < 0 {
		panic("min and max must be positive")
	}
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	if n < 0 {
		return ""
	}
	var sb strings.Builder
	k := len(Alphabet)

	for i := 0; i < n; i++ {
		c := Alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(50)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	return fmt.Sprintf("%v@email.com", RandomString(6))
}

func RandomAccount(owner string) db.Accounts {
	return db.Accounts{
		ID:       RandomInt(1, 1000),
		Owner:    owner,
		Balance:  RandomMoney(),
		Currency: RandomCurrency(),
	}
}

func RandomUser(t *testing.T) (user db.Users, password string) {
	password = RandomString(6)
	hashedPassword, err := HashingPassword(password)
	require.NoError(t, err)
	user = db.Users{
		Username:       RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       RandomOwner(),
		Email:          RandomEmail(),
	}
	return
}
