package external

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	ServerHost               string
	SQSVideoProcessQueueName string
	SQSVideoStatusQueueName  string
	CognitoUserPoolClientId  string
	Region                   string
}

var (
	runOnce sync.Once
	config  Config
)

func GetConfig() Config {
	runOnce.Do(func() {
		cfg, err := initConfig()
		if err != nil {
			fmt.Println(context.Background(), err, "could not load usecase configuration")
		}
		config = Config{
			ServerHost:               cfg.GetString("server.host"),
			SQSVideoProcessQueueName: cfg.GetString("sqs_video_process_queue_name"),
			SQSVideoStatusQueueName:  cfg.GetString("sqs_video_status_queue_name"),
			CognitoUserPoolClientId:  cfg.GetString("COGNITO_CLIENT_USER_POOL_ID"),
			Region:                   cfg.GetString("region"),
		}
	})

	return config
}

func (c Config) IsLocalEnvironment() bool {
	environment := os.Getenv("environment")
	return environment == "development" || environment == "local"
}

func initConfig() (*viper.Viper, error) {
	cfg := viper.New()
	var err error
	initDefaults(cfg)
	// workaround because viper does not resolve envs when unmarshalling
	for _, key := range cfg.AllKeys() {
		val := cfg.Get(key)
		cfg.Set(key, val)
	}
	return cfg, err
}

func initDefaults(config *viper.Viper) {
	config.SetDefault("server.host", "0.0.0.0:8000")
	config.SetDefault("sqs_video_process_queue_name", "video-process-queue")
	config.SetDefault("sqs_video_status_queue_name", "video-status-api-queue")
	config.SetDefault("region", "us-east-1")
	config.SetDefault("COGNITO_CLIENT_USER_POOL_ID", os.Getenv("COGNITO_CLIENT_USER_POOL_ID"))
}
