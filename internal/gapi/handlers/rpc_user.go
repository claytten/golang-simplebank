package gapiHandler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	gapiConverter "github.com/claytten/golang-simplebank/internal/gapi/converter"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *gapiHandlerSetup) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashingPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hashed password %s", err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
		HashedPassword: hashedPassword,
	}
	user, err := s.server.DB.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "User cannot use. Please try again.")
			}
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user : %s", err.Error())
	}

	res := &pb.CreateUserResponse{
		User: gapiConverter.ConvertUser(user),
	}
	return res, nil
}

func (s *gapiHandlerSetup) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.server.DB.GetUserUsingEmail(ctx, req.GetEmail())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "User Not Found")
		}
		return nil, status.Error(codes.Internal, "Cannot Login!. username/email or password is not match.")
	}

	err = util.ComparePassword(user.HashedPassword, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "Cannot Login!. username/email or password is not match.")
	}

	accessToken, accessPayload, err := s.server.Token.CreateToken(user.Email, s.server.Config.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, refreshPayload, err := s.server.Token.CreateToken(user.Email, s.server.Config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	extractMetadata := gapiConverter.ExtractMetadata(ctx, s.server)
	session, err := s.server.DB.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    extractMetadata.UserAgent,
		ClientIp:     extractMetadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
		CreatedAt:    time.Now(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  gapiConverter.ConvertUser(user),
	}
	return res, nil
}

func (s *gapiHandlerSetup) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	authPayload, err := gapiConverter.ConvertToken(req.GetAccessToken(), s.server.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "access denied")
	}

	userHeader, err := s.server.DB.GetUserUsingEmail(ctx, authPayload.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "account doesn't belong to authenticated user")
	}

	if userHeader.Username != req.GetUsername() {
		ff := fmt.Sprintf("%s %s", userHeader.Username, req.GetUsername())
		return nil, status.Error(codes.Unauthenticated, ff)
	}

	user, err := s.server.DB.GetUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.NotFound, "User Not Found")
	}

	return &pb.GetUserResponse{
		User: gapiConverter.ConvertUser(user),
	}, nil
}

func (s *gapiHandlerSetup) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	username, err := gapiConverter.CheckOwnUser(req.GetUsername(), req.GetOldPassword(), req.GetAccessToken(), s.server, ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.server.DB.GetUser(ctx, username)
	if err != nil {
		return nil, status.Error(codes.Internal, "User Not Found")
	}
	if req.GetFullName() == "" {
		req.FullName = user.FullName
	}

	if req.GetEmail() == "" {
		req.Email = user.Email
	}

	updatedUser, err := s.server.DB.UpdateUser(ctx, db.UpdateUserParams{
		Username: user.Username,
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  true,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  true,
		},
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "User Cannot Updated")
	}

	res := &pb.UpdateProfileResponse{
		User: gapiConverter.ConvertUser(updatedUser),
	}

	return res, nil
}

func (s *gapiHandlerSetup) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	username, err := gapiConverter.CheckOwnUser(req.GetUsername(), req.GetOldPassword(), req.GetAccessToken(), s.server, ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.server.DB.GetUser(ctx, username)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User Not Found")
	}

	newPassword, err := util.HashingPassword(req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "Password can't use. Please Try Another Password")
	}

	updatedUser, err := s.server.DB.UpdateUser(ctx, db.UpdateUserParams{
		Username: user.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "User Cannot Updated")
	}

	res := &pb.UpdatePasswordResponse{
		User: gapiConverter.ConvertUser(updatedUser),
	}

	return res, nil
}
