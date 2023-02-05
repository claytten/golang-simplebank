package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/claytten/golang-simplebank/internal/api/token"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type Server struct {
	Engine *gin.Engine
	DB     db.Store
	Config util.Config
	Token  token.Maker
}

func SetupServer(config util.Config, store db.Store) (*Server, error) {
	// can change to JWT or Paseto maker for token
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		DB:     store,
		Config: config,
		Engine: nil,
		Token:  tokenMaker,
	}

	return server, nil
}

// for testing purpose
func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := SetupServer(config, store)
	require.NoError(t, err)

	return server
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.Engine.Run(address)
}
