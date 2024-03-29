// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Accounts, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Accounts, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entries, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Sessions, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfers, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (Users, error)
	DeleteAccount(ctx context.Context, id int64) error
	DeleteAllAccount(ctx context.Context) error
	GetAccount(ctx context.Context, id int64) (Accounts, error)
	GetEntry(ctx context.Context, id int64) (Entries, error)
	GetListsTransfers(ctx context.Context, arg GetListsTransfersParams) ([]Transfers, error)
	GetSession(ctx context.Context, id uuid.UUID) (Sessions, error)
	GetTotalPageListsAccounts(ctx context.Context, owner string) (int64, error)
	GetTotalPageListsTransfers(ctx context.Context) (int64, error)
	GetTotalPageListsTransfersSpesific(ctx context.Context, fromAccountID int64) (int64, error)
	GetTransferByFromAccountId(ctx context.Context, fromAccountID int64) (Transfers, error)
	GetTransferById(ctx context.Context, id int64) (Transfers, error)
	GetTransferByToAccountId(ctx context.Context, toAccountID int64) (Transfers, error)
	GetUser(ctx context.Context, username string) (Users, error)
	GetUserUsingEmail(ctx context.Context, email string) (Users, error)
	ListsAccounts(ctx context.Context, arg ListsAccountsParams) ([]Accounts, error)
	ListsEntries(ctx context.Context, arg ListsEntriesParams) ([]Entries, error)
	ListsTransfers(ctx context.Context, arg ListsTransfersParams) ([]Transfers, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Accounts, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (Users, error)
}

var _ Querier = (*Queries)(nil)
