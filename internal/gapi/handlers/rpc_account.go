package gapiHandler

import (
	"context"
	"database/sql"
	"time"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	gapiConverter "github.com/claytten/golang-simplebank/internal/gapi/converter"
	"github.com/claytten/golang-simplebank/pb"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *gapiHandlerSetup) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	authPayload, err := gapiConverter.AuthorizeUser(ctx, s.server)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	username, err := gapiConverter.CheckOwnUser(req.GetUsername(), req.GetOldPassword(), authPayload.Email, s.server, ctx)
	if err != nil {
		return nil, err
	}

	arg := db.CreateAccountParams{
		Owner:    username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := s.server.DB.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return nil, status.Error(codes.PermissionDenied, err.Error())
			}
		}
		return nil, status.Error(codes.Internal, "Cannot Create Account")
	}

	return &pb.CreateAccountResponse{
		Account: gapiConverter.ConvertAccount(account),
	}, nil
}
func (s *gapiHandlerSetup) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	_, err := gapiConverter.AuthorizeUser(ctx, s.server)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	account, err := s.server.DB.GetAccount(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &pb.GetAccountResponse{
		Account: gapiConverter.ConvertAccount(account),
	}
	return res, nil
}
func (s *gapiHandlerSetup) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	authPayload, err := gapiConverter.AuthorizeUser(ctx, s.server)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	_, err = gapiConverter.CheckOwnUser(req.GetUsername(), req.GetOldPassword(), authPayload.Email, s.server, ctx)
	if err != nil {
		return nil, err
	}

	args := db.UpdateAccountParams{
		ID:        req.GetId(),
		Balance:   req.Balance,
		UpdatedAt: time.Now(),
	}

	account, err := s.server.DB.UpdateAccount(ctx, args)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot update balance account")
	}

	res := &pb.UpdateAccountResponse{
		Account: gapiConverter.ConvertAccount(account),
	}

	return res, nil
}
func (s *gapiHandlerSetup) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	authPayload, err := gapiConverter.AuthorizeUser(ctx, s.server)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	_, err = gapiConverter.CheckOwnUser(req.GetUsername(), req.GetOldPassword(), authPayload.Email, s.server, ctx)
	if err != nil {
		return nil, err
	}
	err = s.server.DB.DeleteAccount(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot delete account")
	}

	res := &pb.DeleteAccountResponse{
		Message: "Delete Successfully",
	}
	return res, nil
}
func (s *gapiHandlerSetup) TransferTxAccount(ctx context.Context, req *pb.TransferTxAccountRequest) (*pb.TransferTxAccountResponse, error) {
	authPayload, err := gapiConverter.AuthorizeUser(ctx, s.server)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	_, err = gapiConverter.CheckOwnUser(req.GetUsername(), req.GetOldPassword(), authPayload.Email, s.server, ctx)
	if err != nil {
		return nil, err
	}

	err = gapiConverter.ValidateAccount(ctx, s.server.DB, req.GetFromAccountID(), req.GetToAccountID(), req.Currency)
	if err != nil {
		return nil, err
	}

	arg := db.TransferTxParams{
		FromAccountID: req.GetFromAccountID(),
		ToAccountID:   req.GetToAccountID(),
		Amount:        req.GetAmount(),
	}

	result, err := s.server.DB.TransferTx(ctx, arg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := gapiConverter.ConvertTransferTx(result)

	return res, nil
}
