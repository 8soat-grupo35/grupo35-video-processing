package external

import (
	"context"
	"fmt"
	"log"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/guregu/dynamo/v2"
)

var (
	DB *dynamo.DB
)

func ConectaDB(config Config) *dynamo.DB {

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Println(err.Error())
		log.Panic("Erro na conexao com banco de dados")
	}

	fmt.Println(config.Environment)
	if config.Environment == "development" {
		baseURL := "http://localstack:4566"
		cfg.BaseEndpoint = &baseURL
	}

	DB = dynamo.New(cfg)

	err = DB.CreateTable("production_order", entities.Video{}).OnDemand(true).Run(context.TODO())

	if err != nil {
		log.Println(err.Error())
	}

	return DB
}
