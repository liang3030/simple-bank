package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liang3030/simple-bank/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBear   = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeaser := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeaser) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeaser)
		if len(fields) < 2 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBear {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()

	}
}
