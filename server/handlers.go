package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func statusPage(c echo.Context) error {
	return c.String(http.StatusOK, "api is online!")
}
