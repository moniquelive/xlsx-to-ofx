// Package router sets up web routes
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/moniquelive/xlsx-to-ofx/handler"
)

// SetupRoutes defines the app web routes
func SetupRoutes(app *fiber.App) {
	app.Static("/css", "./web/css").
		Static("/images", "./web/images").
		Static("/js", "./web/js")

	app.Get("/status", monitor.New()).
		Get("/", handler.CsrfProtection, handler.Index)

	app.Post("/convert", handler.CsrfProtection, handler.DoConvert)
}
