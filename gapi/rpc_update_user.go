package gapi

import (
	"context"
	"database/sql"
	"log"
	"time"

	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/pb"
	"github.com/liang3030/simple-bank/util"
	"github.com/liang3030/simple-bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, InvalidArgumentError(violations)
	}

	// It is used to extract metadata from the context for login user.
	// Put it here to show how to use it because login is not implemented yet.
	mdt := server.extractMetadata(ctx)
	log.Printf("metadata: %v", mdt)

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{String: req.GetFullName(), Valid: req.FullName != nil},
		Email:    sql.NullString{String: req.GetEmail(), Valid: req.Email != nil},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())

		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}
		arg.HashedPassword = sql.NullString{String: hashedPassword, Valid: true}
		arg.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldValidation("username", err))
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldValidation("password", err))
		}
	}

	if req.FullName != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldValidation("password", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateFullname(req.GetFullName()); err != nil {
			violations = append(violations, fieldValidation("full_name", err))
		}
	}
	return violations
}
