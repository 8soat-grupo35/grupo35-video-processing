package server

import (
	"context"
	"fmt"

	"net/http"

	_ "github.com/8soat-grupo35/grupo35-video-processing/docs"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/adapters/wrappers"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/api/handlers"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/external"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/gateways"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/usecases"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
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
	database := external.ConectaDB(cfg)

	awsConfig, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("AWS Config loaded")

	app := echo.New()
	app.GET("/swagger/*", echoSwagger.WrapHandler)

	app.GET("/", func(echo echo.Context) error {
		return echo.JSON(http.StatusOK, "Alive")
	})

	sqsClient := wrappers.NewSQSClient(awsConfig)
	s3Client := wrappers.NewS3Client(awsConfig)

	videoHandler := handlers.NewVideoHandler(
		usecases.NewTransferVideo(
			gateways.NewS3Gateway(s3Client),
		),
		usecases.NewVideoProcessSender(
			gateways.NewSQSGateway(sqsClient, cfg.SQSVideoProcessQueueName, 10),
		),
		usecases.NewVideoRepository(
			gateways.NewDynamoGateway(external.NewDynamoAdapter(database)),
		),
	)
	app.Use(AuthenticationMiddleware())
	{
		videoGroupV1 := app.Group("/v1/videos")
		videoGroupV1.POST("/upload", videoHandler.Upload)
	}

	return app
}
