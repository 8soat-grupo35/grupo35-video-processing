package server

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/8soat-grupo35/fastfood-order-production/docs"
	"github.com/8soat-grupo35/fastfood-order-production/internal/external"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Start(cfg external.Config) {
	fmt.Println(context.Background(), fmt.Sprintf("Starting a server at http://%s", cfg.ServerHost))
	app := newApp(cfg)
	app.Logger.Fatal(app.Start(cfg.ServerHost))
}

// @title Swagger Fastfood App API
// @version 1.0
// @description This is a sample API from Fastfood App.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /v1
func newApp(cfg external.Config) *echo.Echo {
	external.ConectaDB(cfg.DatabaseConfig.Host, cfg.DatabaseConfig.User, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.DbName, cfg.DatabaseConfig.Port)

	app := echo.New()
	app.GET("/swagger/*", echoSwagger.WrapHandler)
	app.GET("/", func(echo echo.Context) error {
		return echo.JSON(http.StatusOK, "Alive")
	})

	return app
}
