package gapiConverter

import (
	"context"
	"database/sql"

	"github.com/claytten/golang-simplebank/internal/api/token"
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

func ConvertToken(token string, maker token.Maker) (*token.Payload, error) {
	authPayload, err := maker.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	return authPayload, nil
}

func CheckOwnUser(username, oldPassword, token string, server *gapi.Server, ctx context.Context) (string, error) {
	authPayload, err := ConvertToken(token, server.Token)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "access denied")
	}

	//finding user by email
	userHeader, err := server.DB.GetUserUsingEmail(ctx, authPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", status.Error(codes.NotFound, "User Not Found and Not Authorized")
		}

		return "", status.Error(codes.Internal, "account doesn't belong to authenticated user")
	}

	// checking if username is provided at header
	if username != userHeader.Username {
		return "", status.Error(codes.Internal, "User Status Unauthorized")
	}

	// checking if user typing old password and new password is same
	err = util.ComparePassword(userHeader.HashedPassword, oldPassword)
	if err != nil {
		return "", status.Error(codes.Internal, "new password / old password doesn't match")
	}

	return userHeader.Username, nil
}
