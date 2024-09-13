package http

import (
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"net/http"
	"time"
)

type Response struct {
	Version string    `json:"version"`
	Name    string    `json:"name"`
	Uptime  time.Time `json:"uptime"`
}

type PingHandler struct {
	config config.Config
}

func NewPingHandler(cfg config.Config) *PingHandler {
	return &PingHandler{
		config: cfg,
	}
}

func (s *PingHandler) Ping(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Response{
		Version: s.config.ProjectVersion,
		Name:    s.config.ProjectName,
		Uptime:  time.Now().UTC(),
	})
}
