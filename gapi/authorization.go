package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/liang3030/simple-bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "Bearer"
)

func (server *Server) authorizedUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	// <authType> <token>
	// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}

	authType := fields[0]
	if authType != authorizationType {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)

	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	return payload, nil
}
