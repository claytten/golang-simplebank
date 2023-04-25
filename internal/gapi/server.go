package gapi

import (
	"fmt"

	"github.com/claytten/golang-simplebank/internal/api/token"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/internal/worker"
)

type Server struct {
	DB             db.Store
	Config         util.Config
	Token          token.Maker
	TaskDistrbutor worker.TaskDistributor
}

func SetupServer(config util.Config, store db.Store, taskDistrbutor worker.TaskDistributor) (*Server, error) {
	// can change to JWT or Paseto maker for token
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		DB:             store,
		Config:         config,
		Token:          tokenMaker,
		TaskDistrbutor: taskDistrbutor,
	}

	return server, nil
}
