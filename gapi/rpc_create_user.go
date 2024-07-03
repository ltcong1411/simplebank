package gapi

import (
	"context"
	"fmt"
	"time"

	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/val"
	"simplebank/worker"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// TODO: use db transaction
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task to send verify email: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (validations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		validations = append(validations, fieldViolation("username", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		validations = append(validations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		validations = append(validations, fieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		validations = append(validations, fieldViolation("password", err))
	}

	return validations
}
