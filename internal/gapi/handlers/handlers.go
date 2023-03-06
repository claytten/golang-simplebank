package gapiHandler

import (
	"github.com/claytten/golang-simplebank/internal/gapi"
	"github.com/claytten/golang-simplebank/pb"
)

type gapiHandlerSetup struct {
	pb.UnimplementedSimplebankServer
	server *gapi.Server
}

func NewGapiHandlerSetup(server *gapi.Server) *gapiHandlerSetup {
	return &gapiHandlerSetup{server: server}
}
