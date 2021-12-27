package server

import (
	"context"

	pb "github.com/ip75/meteostation/proto/api"
	"github.com/ip75/meteostation/storage"
	log "github.com/sirupsen/logrus"
)

type MeteostationServer struct {
	pb.UnimplementedMeteostationServiceServer
}

func NewMeteostationServer() *MeteostationServer {
	return &MeteostationServer{}
}

func (s *MeteostationServer) GetMeteoData(ctx context.Context, filter *pb.Filter) (*pb.MeteoData, error) {
	log.Printf("GetMeteoData:\n\tfrom: %s\n\tto: %s\n\tgranularity: %d", filter.From, filter.To, filter.Granularity)

	data, err := storage.PG.GetMeteoData(filter)
	if err != nil {
		log.WithFields(log.Fields{
			"from": filter.From,
			"to":   filter.To,
		}).Error(err, "Error when fetching meteodata")
		return nil, err
	}

	log.Printf("GetMeteoData: fetched %d records with meteodata", data.TotalCount)
	return data, err
}
