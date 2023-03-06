package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/gapi"
	gapiHandlerSetup "github.com/claytten/golang-simplebank/internal/gapi/handlers"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig("app", ".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB : ", err)
	}

	store := db.NewStore(conn)
	RunGrpcServer(config, store)
}

func RunGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.SetupServer(config, store)
	if err != nil {
		log.Fatal("cannot create gRPC Server : ", err)
	}

	grpcServer := grpc.NewServer()
	handlers := gapiHandlerSetup.NewGapiHandlerSetup(server)
	pb.RegisterSimplebankServer(grpcServer, handlers)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot listen TCP gRPC : ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC Server : ", err)
	}
}

func RunGinServer(config util.Config, store db.Store) {
	server, err := api.SetupServer(config, store)
	if err != nil {
		log.Fatal("cannot create HTTP server : ", err)
	}
	routes.ApplyAllPublicRoutes(server)

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start HTTP server : ", err)
	}
}
