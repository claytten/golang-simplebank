package gapi

import (
	"fmt"

	"github.com/claytten/golang-simplebank/internal/api/token"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
)

type Server struct {
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
		Token:  tokenMaker,
	}

	return server, nil
}
