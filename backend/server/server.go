package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ip75/meteostation/config"
	gw "github.com/ip75/meteostation/proto/api"
)

var (
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-meteostaion-server-endpoint", fmt.Sprintf("localhost:%d", config.C.General.GrpcServicePort), "Meteostation gRPC server endpoint")
)

// start web hosting of static
func StartHosting(ctx context.Context) error {
	fs := http.FileServer(http.Dir(config.C.General.WebStaticDir))
	http.Handle("/", fs)

	go func() {
		<-ctx.Done()
		// put code to release resources when everythiing is stopping
	}()

	log.Printf("start client hosting from %s at port %d ...", config.C.General.WebStaticDir, config.C.General.HttpPort)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.C.General.HttpPort), nil)
}

// Register gRPC server endpoint
// Note: Make sure the gRPC server is running properly and accessible
func StartGrpcService(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.C.General.GrpcServicePort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	gw.RegisterMeteostationServiceServer(grpcServer, NewMeteostationServer())

	go func() {
		defer grpcServer.GracefulStop()
		<-ctx.Done()
	}()

	log.Printf("start grpc server at port %d ...", config.C.General.GrpcServicePort)

	return grpcServer.Serve(listener)
}

func StartGateway(ctx context.Context) error {

	// register gRPC client handler to meteostation server
	mux := runtime.NewServeMux()
	if err := gw.RegisterMeteostationServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
		return err
	}

	log.Printf("start http gateway at port %d ...", config.C.General.GrpcHttpGatewayPort)

	// Start HTTP server (and proxy calls to gRPC server endpoint) through mux
	return http.ListenAndServe(fmt.Sprintf(":%d", config.C.General.GrpcHttpGatewayPort), mux)
}
