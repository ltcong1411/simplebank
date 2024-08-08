package gapi

import (
	"context"
	"errors"
	"time"

	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/val"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != util.BankerRole && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	res := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (validations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		validations = append(validations, fieldViolation("username", err))
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			validations = append(validations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			validations = append(validations, fieldViolation("email", err))
		}
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			validations = append(validations, fieldViolation("password", err))
		}
	}

	return validations
}
