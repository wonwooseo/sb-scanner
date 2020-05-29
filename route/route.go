package route

import (
	"sb-scanner/api"

	"github.com/labstack/echo/v4"
)

// AddRoutes _
func AddRoutes(e *echo.Echo) {
	e.GET("/latest", api.GetLatestHandler)
}
