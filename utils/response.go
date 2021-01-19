package utils

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func ResponseByContentType(context echo.Context, status int, data interface{}) error {
	acceptHeader := context.Request().Header.Get("Accept")
	if strings.Contains(strings.ToLower(acceptHeader), "xml") {
		return context.XML(status, data)
	} else {
		return context.JSON(status, data)
	}
}
