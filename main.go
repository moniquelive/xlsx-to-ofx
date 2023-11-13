// package comment
package main

import (
	"embed"
	"net/http"
	"os"
	"runtime"

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
)

//go:embed web
var embedFS embed.FS

func main() {
	engine := getEngine(runtime.GOOS)
	app := setupFiberApp(engine)

	router.SetupRoutes(app)
	app.All("/*", filesystem.New(filesystem.Config{PathPrefix: "web", Root: engine.FileSystem}))

	port := getPort()
	log.Fatal(app.Listen(port))
}

func getEngine(goos string) *html.Engine {
	var engine *html.Engine
	if goos == "linux" {
		engine = html.NewFileSystem(http.FS(embedFS), ".gohtml")
		engine.Directory = "web"
	} else {
		engine = html.NewFileSystem(http.FS(os.DirFS("web")), ".gohtml")
		engine.Reload(true)
		engine.Debug(true)
	}
	return engine
}

func setupFiberApp(engine *html.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})
	app.Use(favicon.New()).
		Use(logger.New()).
		Use(helmet.New()).
		Use(recover.New()).
		Use(cors.New(cors.Config{
			AllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
			AllowHeaders: "Origin, Content-Type, Accept",
		}))
	return app
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "localhost:9090" // default port
	}
	return port
}
