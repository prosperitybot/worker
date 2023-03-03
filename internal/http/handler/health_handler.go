package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"go.uber.org/zap"
)

type HealthHandler struct {
	Db *sqlx.DB
}

func (h HealthHandler) GETHealth(c echo.Context) error {

	if err := h.Db.PingContext(c.Request().Context()); err != nil {
		logger.Error(c.Request().Context(), "Error pinging database", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Error pinging database")
	}

	return c.String(http.StatusOK, "OK")
}
