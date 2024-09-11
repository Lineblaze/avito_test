package httpServer

import (
	"fmt"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"zadanie-6105/backend/internal/common/config"
	"zadanie-6105/backend/pkg/httpErrorHandler"
	"zadanie-6105/backend/pkg/logger"
)

// Server struct
type Server struct {
	fiber     *fiber.App
	cfg       *config.Config
	apiLogger *logger.ApiLogger
}

func NewServer(cfg *config.Config, apiLogger *logger.ApiLogger, handler *httpErrorHandler.HttpErrorHandler) *Server {
	return &Server{
		fiber: fiber.New(fiber.Config{
			ErrorHandler: handler.Handler,
			JSONEncoder:  gojson.Marshal,
			JSONDecoder:  gojson.Unmarshal,
		}),
		cfg:       cfg,
		apiLogger: apiLogger,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.fiber, s.apiLogger); err != nil {
		s.apiLogger.Fatalf("Cannot map handlers: ", err)
	}
	s.apiLogger.Infof("Start server on port: %s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	if err := s.fiber.Listen(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)); err != nil {
		s.apiLogger.Fatalf("Error starting Server: ", err)
	}
	return nil
}
