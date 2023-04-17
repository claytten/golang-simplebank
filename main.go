package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/claytten/golang-simplebank/doc/statik"
	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/gapi"
	gapiHandlerSetup "github.com/claytten/golang-simplebank/internal/gapi/handlers"
	gapiLogger "github.com/claytten/golang-simplebank/internal/gapi/logger"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/claytten/golang-simplebank/pb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// checking file app.env is exists
	config, err := util.LoadConfig("app", ".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		runLogFile, err := RunLogProduction("myapp", "log")
		if err != nil {
			log.Fatal().Err(err).Msg("cannot create log file")
		}
		multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)

		log.Logger = zerolog.New(multi)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to DB")
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
		log.Fatal().Err(err).Msg("cannot create gRPC Server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapiLogger.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	handlers := gapiHandlerSetup.NewGapiHandlerSetup(server)
	pb.RegisterSimplebankServer(grpcServer, handlers)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot listen TCP gRPC")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC Server")
	}
}

func RunGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.SetupServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC Server")
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
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// adding swagger
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik file system")
	}

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFs)))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handlerLogger := gapiLogger.HttpLogger(mux)
	err = http.Serve(listener, handlerLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

func RunGinServer(config util.Config, store db.Store) {
	server, err := api.SetupServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create HTTP server")
	}
	routes.ApplyAllPublicRoutes(server)

	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP server")
	}
}

func RunDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run migration")
	}

	log.Info().Msg("migration successful")
}

func RunLogProduction(logName, folder string) (*os.File, error) {
	// checking file log/myapp.log is exists
	dayString := time.Now().Format("20060102")
	fileName := logName + dayString + ".log"

	if err := util.ValidateFolderAndFile(folder, fileName); err != nil {
		return nil, err
	}

	runLogFile, err := os.OpenFile(
		folder+"/"+fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	return runLogFile, err
}
