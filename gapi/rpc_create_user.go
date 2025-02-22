package gapi

import (
	"context"
	"log"

	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/pb"
	"github.com/liang3030/simple-bank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// It is used to extract metadata from the context for login user.
	// Put it here to show how to use it because login is not implemented yet.
	mdt := server.extractMetadata(ctx)
	log.Printf("metadata: %v", mdt)

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user name already existed")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}
