package gapi

import (
	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/pb"
	"github.com/liang3030/simple-bank/token"
	"github.com/liang3030/simple-bank/util"
)

// Server GRPC requests for banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.IStore
	tokenMaker token.Maker
}

// NewServer creates a new GRPC server.
func NewServer(config util.Config, store db.IStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker}

	return server, nil
}
