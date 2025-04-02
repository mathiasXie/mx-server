package server

import (
	"context"

	"github.com/mathiasXie/gin-web/applications/function-rpc/internal/weather"
	"github.com/mathiasXie/gin-web/applications/function-rpc/proto/pb/proto"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

type Server struct {
	proto.UnimplementedFunctionServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetWeather(ctx context.Context, req *proto.GetWeatherRequest) (*proto.GetWeatherResponse, error) {
	weather, err := weather.GetWeatherByCityName(req.Location)
	if err != nil {
		return nil, err
	}
	return &proto.GetWeatherResponse{
		Temp:      weather.Now.Temp,
		FeelsLike: weather.Now.FeelsLike,
		Icon:      weather.Now.Icon,
		Text:      weather.Now.Text,
		Wind360:   weather.Now.Wind360,
		WindDir:   weather.Now.WindDir,
		WindScale: weather.Now.WindScale,
		WindSpeed: weather.Now.WindSpeed,
		Humidity:  weather.Now.Humidity,
		Precip:    weather.Now.Precip,
		Pressure:  weather.Now.Pressure,
		Vis:       weather.Now.Vis,
		Cloud:     weather.Now.Cloud,
		Dew:       weather.Now.Dew,
	}, nil
}

func (s *Server) GetWeatherReport(ctx context.Context, req *proto.GetWeatherReportRequest) (*proto.GetWeatherReportResponse, error) {
	report, err := weather.GetWeatherReport(req.Location, req.Lang)

	logger.CtxInfo(ctx, "GetWeatherReport", "report", report)
	if err != nil {
		return nil, err
	}
	return &proto.GetWeatherReportResponse{Report: report}, nil
}
