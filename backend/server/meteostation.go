package server

import (
	"context"

	pb "github.com/ip75/meteostation/proto/api"
	"github.com/ip75/meteostation/storage"
)

type MeteostationServer struct {
	pb.UnimplementedMeteostationServiceServer
}

func NewMeteostationServer() MeteostationServer {
	return MeteostationServer{}
}

func (s *MeteostationServer) GetMeteoData(ctx context.Context, filter *pb.Filter) (*pb.MeteoData, error) {
	return storage.PG.GetMeteoData(filter)
}
