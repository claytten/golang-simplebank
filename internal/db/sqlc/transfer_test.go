package db_test

import (
	"context"
	"math"
	"testing"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

func TestGetListsTransfers(t *testing.T) {
	fromAccount, _, amount := CreateTransferTx(t, 5)

	arg := db.GetListsTransfersParams{
		FromAccountID: fromAccount.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.GetListsTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	var totalAmount int64
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		totalAmount += transfer.Amount
	}
	require.Equal(t, amount, totalAmount)
}

func TestGetTotalPageListsTransfersSpesific(t *testing.T) {
	total := int(util.RandomInt(1, 10))
	fromAccount, _, _ := CreateTransferTx(t, total)

	transfers, err := testQueries.GetTotalPageListsTransfersSpesific(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Equal(t, int64(total), transfers)
}

func TestGetTotalPageListsTransfers(t *testing.T) {
	limit := 5
	//if limit 5 and create transfer just 5
	//page == 1
	CreateTransferTx(t, limit)

	transfers, err := testQueries.GetTotalPageListsTransfers(context.Background())
	total_page := math.Ceil(float64(transfers) / float64(5))
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.GreaterOrEqual(t, total_page, float64(1))
}

func TestGetTransferFromAccountId(t *testing.T) {
	fromAccount, _, _ := CreateTransferTx(t, 1)

	transfer, err := testQueries.GetTransferByFromAccountId(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, fromAccount.ID, transfer.FromAccountID)
	require.NotZero(t, transfer.ID)
}

func TestGetTransferToAccountId(t *testing.T) {
	_, toAccount, _ := CreateTransferTx(t, 1)

	transfer, err := testQueries.GetTransferByToAccountId(context.Background(), toAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, toAccount.ID, transfer.ToAccountID)
	require.NotZero(t, transfer.ID)
}

func TestListsTransfers(t *testing.T) {
	total := 5
	CreateTransferTx(t, total)

	arg := db.ListsTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListsTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.GreaterOrEqual(t, total, len(transfers))
}
