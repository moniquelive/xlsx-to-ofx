// package comment
package main

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/moniquelive/xlsx-to-ofx/router"
	"net/http"
	"os"
	"runtime"
)

//go:embed web
var embedFS embed.FS

func main() {
	var engine *html.Engine
	if runtime.GOOS == "linux" {
		engine = productionEngine()
	} else {
		engine = developmentEngine()
	}
	fiberConfig := fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	}
	corsConfig := cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}
	app := fiber.New(fiberConfig)
	app.Use(favicon.New()).
		Use(logger.New()).
		Use(helmet.New()).
		Use(recover.New()).
		Use(cors.New(corsConfig))

	router.SetupRoutes(app)

	if runtime.GOOS == "linux" {
		// embedded static stuff
		app.Get("/*", filesystem.New(filesystem.Config{PathPrefix: "web", Root: http.FS(embedFS)}))
		log.Fatal(app.Listen(":9090"))
	} else {
		log.Fatal(app.Listen("localhost:9090"))
	}
}

func developmentEngine() (engine *html.Engine) {
	engine = html.NewFileSystem(http.FS(os.DirFS("web")), ".gohtml")
	engine.Reload(true)
	engine.Debug(true)
	return
}

func productionEngine() (engine *html.Engine) {
	engine = html.NewFileSystem(http.FS(embedFS), ".gohtml")
	engine.Directory = "web"
	return
}
