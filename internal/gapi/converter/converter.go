package gapiConverter

import (
	"context"
	"database/sql"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/gapi"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user db.Users) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		UpdatedAt:         timestamppb.New(user.UpdatedAt),
	}
}

func ConvertAccount(account db.Accounts) *pb.Account {
	return &pb.Account{
		Id:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: timestamppb.New(account.CreatedAt),
		UpdatedAt: timestamppb.New(account.UpdatedAt),
	}
}

func ConvertTransfer(transfer db.Transfers) *pb.Transfer {
	return &pb.Transfer{
		Id:            transfer.ID,
		FromAccountId: transfer.FromAccountID,
		ToAccountId:   transfer.ToAccountID,
		Amount:        transfer.Amount,
		CreatedAt:     timestamppb.New(transfer.CreatedAt),
		UpdatedAt:     timestamppb.New(transfer.UpdatedAt),
	}
}

func ConvertEntry(entry db.Entries) *pb.Entries {
	return &pb.Entries{
		Id:        entry.ID,
		AccountId: entry.AccountID,
		Amount:    entry.Amount,
		CreatedAt: timestamppb.New(entry.CreatedAt),
		UpdatedAt: timestamppb.New(entry.UpdatedAt),
	}
}

func ConvertTransferTx(transfer db.TransferTxResult) *pb.TransferTxAccountResponse {
	return &pb.TransferTxAccountResponse{
		Transfer:    ConvertTransfer(transfer.Transfer),
		FromAccount: ConvertAccount(transfer.FromAccount),
		ToAccount:   ConvertAccount(transfer.ToAccount),
		FromEntry:   ConvertEntry(transfer.FromEntry),
		ToEntry:     ConvertEntry(transfer.ToEntry),
	}
}

func CheckOwnUser(username, oldPassword, email string, server *gapi.Server, ctx context.Context) (string, error) {
	//finding user by email
	userHeader, err := server.DB.GetUserUsingEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", status.Error(codes.NotFound, "User Not Found and Not Authorized")
		}

		return "", status.Error(codes.Internal, "account doesn't belong to authenticated user")
	}

	// checking if username is provided at header
	if username != userHeader.Username {
		return "", status.Error(codes.PermissionDenied, "User Status Permission Denied")
	}

	// checking if user typing old password and new password is same
	err = util.ComparePassword(userHeader.HashedPassword, oldPassword)
	if err != nil {
		return "", status.Error(codes.Internal, "new password / old password doesn't match")
	}

	return userHeader.Username, nil
}

func ValidateAccount(ctx context.Context, db db.Store, from_account_id, to_account_id int64, currency string) error {
	fromAccount, err := db.GetAccount(ctx, from_account_id)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	toAccount, err := db.GetAccount(ctx, from_account_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return status.Error(codes.NotFound, err.Error())
		}

		return status.Error(codes.Internal, err.Error())
	}

	if fromAccount.Currency != currency || toAccount.Currency != currency {
		return status.Error(codes.Internal, "From/To Account Mismatch Currency")
	}
	return nil
}
