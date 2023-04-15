package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/claytten/golang-simplebank/doc/statik"
	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/gapi"
	gapiHandlerSetup "github.com/claytten/golang-simplebank/internal/gapi/handlers"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
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

	RunDBMigration(config.MigrationURL, config.DBSource)

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

	// adding swagger
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik file system : ", err.Error())
	}

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFs)))

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

func RunDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot run migration : ", err)
	}

	log.Println("migration successful")
}
