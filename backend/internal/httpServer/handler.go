package httpServer

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	serverLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"zadanie-6105/backend/internal/delivery/http"
	repository "zadanie-6105/backend/internal/repository"
	useCase "zadanie-6105/backend/internal/usecase"
	"zadanie-6105/backend/pkg/logger"
	storage "zadanie-6105/backend/pkg/storage/postgres"
)

func (s *Server) MapHandlers(app *fiber.App, logger *logger.ApiLogger) error {
	db, err := storage.InitPsqlDB(s.cfg)
	if err != nil {
		return err
	}

	repo := repository.NewPostgresRepository(db)
	useCase := useCase.NewUseCase(repo)
	handler := http.NewHandler(useCase, logger)

	app.Use(serverLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{},
		AllowHeaders: []string{},
	}))

	group := app.Group("api")
	http.MapRoutes(group, handler)

	return nil
}
