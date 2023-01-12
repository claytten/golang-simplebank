package db_test

import (
	"context"
	"database/sql"
	"math"
	"testing"
	"time"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

/** start testing normally **/
func CreateRandomAccount(t *testing.T) db.Accounts {
	user := CreateRandomUser(t)

	arg := db.CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.WithinDuration(t, account1.UpdatedAt, account2.UpdatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := db.UpdateAccountParams{
		ID:        account1.ID,
		Balance:   account1.Balance,
		UpdatedAt: time.Now(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.WithinDuration(t, arg.UpdatedAt, account2.UpdatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListsAccounts(t *testing.T) {
	var lastAccount db.Accounts
	for i := 0; i < 10; i++ {
		lastAccount = CreateRandomAccount(t)
	}

	arg := db.ListsAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListsAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, v := range accounts {
		require.NotEmpty(t, v)
		require.Equal(t, lastAccount.Owner, v.Owner)
	}
}

func TestTotalAccounts(t *testing.T) {
	var lastAccount db.Accounts
	err := testQueries.DeleteAllAccount(context.Background())
	require.NoError(t, err)

	for i := 0; i < 11; i++ {
		lastAccount = CreateRandomAccount(t)
	}
	rowPerPage := 5

	accounts, err := testQueries.GetTotalPageListsAccounts(context.Background(), lastAccount.Owner)
	totalPage := math.Ceil(float64(accounts) / float64(rowPerPage))
	require.NoError(t, err)
	require.Equal(t, totalPage, float64(1))
}

/** end testing normally **/
