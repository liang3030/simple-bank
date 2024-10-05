package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/util"
	_ "github.com/lib/pq" // Make sure the postgres enginee driver is imported
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.IStore) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
