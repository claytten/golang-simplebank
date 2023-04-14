package gapiValidate

import (
	"fmt"

	gapiError "github.com/claytten/golang-simplebank/internal/gapi/error"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ValidateAuthorizeAccountRequest(username, password string) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(username); err != nil {
		violations = append(violations, gapiError.FieldViolation("account_id", err))
	}

	if err := util.ValidatePassword(password); err != nil {
		violations = append(violations, gapiError.FieldViolation("password", err))
	}

	return violations
}

func ValidateCreateAccountRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateAuthorizeAccountRequest(req.GetUsername(), req.GetOldPassword()); err != nil {
		violations = append(violations, err...)
	}

	if ok := util.IsSupportCurrency(req.GetCurrency()); !ok {
		err := fmt.Errorf("%s", req.GetCurrency())
		violations = append(violations, gapiError.FieldViolation("currency not supported", err))
	}

	return violations
}

func ValidateUpdateAccountRequest(req *pb.UpdateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateAuthorizeAccountRequest(req.GetUsername(), req.GetOldPassword()); err != nil {
		violations = append(violations, err...)
	}

	if err := util.ValidateBalance(req.GetBalance()); err != nil {
		violations = append(violations, gapiError.FieldViolation("balance ", err))
	}

	return violations
}

func ValidateTransactionAccountRequest(req *pb.TransferTxAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateAuthorizeAccountRequest(req.GetUsername(), req.GetOldPassword()); err != nil {
		violations = append(violations, err...)
	}

	if ok := util.IsSupportCurrency(req.GetCurrency()); !ok {
		err := fmt.Errorf("%s", req.GetCurrency())
		violations = append(violations, gapiError.FieldViolation("currency not supported", err))
	}

	if err := util.ValidateBalance(req.GetAmount()); err != nil {
		violations = append(violations, gapiError.FieldViolation("balance ", err))
	}

	return violations
}
