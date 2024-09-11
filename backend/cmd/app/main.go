package main

import (
	"log"
	"zadanie-6105/backend/internal/common/config"
	"zadanie-6105/backend/internal/httpServer"
	"zadanie-6105/backend/pkg/httpErrorHandler"
	"zadanie-6105/backend/pkg/logger"
)

func main() {
	log.Println("Starting server")
	v, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Cannot cload config: ", err.Error())
	}
	cfg, err := config.ParseConfig(v)
	if err != nil {
		log.Fatalf("Config parse error", err.Error())
	}
	log.Println("Config loaded")

	appLogger := logger.NewApiLogger(cfg)
	err = appLogger.InitLogger()
	if err != nil {
		log.Fatalf("Cannot init logger: %v", err.Error())
	}
	log.Println("Logger loaded")

	errorHandler := httpErrorHandler.NewErrorHandler(cfg)
	s := httpServer.NewServer(cfg, appLogger, errorHandler)
	if err = s.Run(); err != nil {
		appLogger.Errorf("run")
	}
}
