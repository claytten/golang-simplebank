package gapiValidate

import (
	gapiError "github.com/claytten/golang-simplebank/internal/gapi/error"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, gapiError.FieldViolation("username", err))
	}

	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, gapiError.FieldViolation("password", err))
	}

	if err := util.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, gapiError.FieldViolation("full_name", err))
	}

	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, gapiError.FieldViolation("email", err))
	}

	return violations
}

func ValidateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, gapiError.FieldViolation("username", err))
	}

	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, gapiError.FieldViolation("password", err))
	}

	return violations
}

func ValidateUpdateProfileRequest(req *pb.UpdateProfileRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, gapiError.FieldViolation("full_name", err))
	}

	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, gapiError.FieldViolation("email", err))
	}

	if err := util.ValidatePassword(req.GetOldPassword()); err != nil {
		violations = append(violations, gapiError.FieldViolation("old password", err))
	}

	return violations
}

func ValidateUpdatePasswordRequest(req *pb.UpdatePasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidatePassword(req.GetOldPassword()); err != nil {
		violations = append(violations, gapiError.FieldViolation("old_password", err))
	}

	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, gapiError.FieldViolation("new_password", err))
	}

	return violations
}
