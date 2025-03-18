package gapi

import (
	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/pb"
	"github.com/liang3030/simple-bank/token"
	"github.com/liang3030/simple-bank/util"
	"github.com/liang3030/simple-bank/worker"
)

// Server GRPC requests for banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.IStore
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new GRPC server.
func NewServer(config util.Config, store db.IStore, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:           store,
		config:          config,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
