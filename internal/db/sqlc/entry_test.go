package db_test

import (
	"context"
	"testing"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T, amount int64) *db.Entries {
	account := CreateRandomAccount(t)

	arg := db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.UpdatedAt)

	return &entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t, util.RandomInt(1, 10))
}

func TestGetEntry(t *testing.T) {
	entry := CreateRandomEntry(t, util.RandomInt(1, 10))
	require.NotEmpty(t, entry)

	getEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getEntry)
	require.Equal(t, entry.ID, getEntry.ID)
	require.Equal(t, entry.AccountID, getEntry.AccountID)
	require.Equal(t, entry.Amount, getEntry.Amount)
	require.Equal(t, entry.CreatedAt, getEntry.CreatedAt)
	require.Equal(t, entry.UpdatedAt, getEntry.UpdatedAt)
}

func TestListsEntries(t *testing.T) {
	var lastEntry *db.Entries
	for i := 0; i < 11; i++ {
		lastEntry = CreateRandomEntry(t, util.RandomInt(1, 10))
	}

	arg := db.ListsEntriesParams{
		AccountID: lastEntry.AccountID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListsEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, lastEntry.AccountID, entry.AccountID)
	}
}
