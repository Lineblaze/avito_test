package httpErrorHandler

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"runtime/debug"
	"zadanie-6105/backend/internal/common/config"
)

type HttpErrorHandler struct {
	showUnknownErrors bool
	logger            *slog.Logger
}

func NewErrorHandler(c *config.Config) *HttpErrorHandler {
	return &HttpErrorHandler{
		showUnknownErrors: c.Server.ShowUnknownErrorsInResponse,
	}
}

type responseMsg struct {
	Message string `json:"message"`
}

func (handler *HttpErrorHandler) Handler(c fiber.Ctx, err error) error {
	var response responseMsg
	var statusCode int

	if statusCode == 0 {
		statusCode = fiber.StatusInternalServerError
	}
	if response.Message == "" && handler.showUnknownErrors {
		response.Message = fmt.Sprintf("Error: \n\n %s", err.Error())
	} else if response.Message == "" {
		handler.logger.Error("empty response message", slog.String("error", err.Error()))
		response.Message = "unknown error"
	}

	handler.logger.Error(err.Error(), slog.String("stack", string(debug.Stack())))
	handler.logger.Error(err.Error(), slog.String("error", err.Error()))

	return c.Status(statusCode).JSON(response)
}

func (handler *HttpErrorHandler) StackTraceHandler(c fiber.Ctx, e interface{}) {
	if e == nil {
		e = "contact support"
	}
	handler.logger.Error("",
		slog.String("path", c.Path()),
		slog.String("method", c.Method()),
		slog.Any("error", e),
	)
	c.Status(500).JSON(map[string]interface{}{
		"description": e,
		"status":      "500",
	})
}
