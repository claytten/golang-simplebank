package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/gapi"
	gapiHandlerSetup "github.com/claytten/golang-simplebank/internal/gapi/handlers"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	// RunGinServer(config, store)
	go RunGatewayServer(config, store)
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
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start gRPC Server : ", err)
	}
}

func RunGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.SetupServer(config, store)
	if err != nil {
		log.Fatal("cannot create gRPC Server : ", err)
	}

	// for snackcase
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	handlers := gapiHandlerSetup.NewGapiHandlerSetup(server)
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimplebankHandlerServer(ctx, grpcMux, handlers)
	if err != nil {
		log.Fatal("cannot register handler server : ", err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener : ", err.Error())
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server : ", err.Error())
	}
}

func RunGinServer(config util.Config, store db.Store) {
	server, err := api.SetupServer(config, store)
	if err != nil {
		log.Fatal("cannot create HTTP server : ", err)
	}
	routes.ApplyAllPublicRoutes(server)

	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("cannot start HTTP server : ", err)
	}
}
