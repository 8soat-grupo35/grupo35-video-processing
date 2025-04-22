package external

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	ServerHost               string
	Environment              string
	SQSVideoProcessQueueName string
	SQSVideoStatusQueueName  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
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
			Environment:              cfg.GetString("environment"),
			SQSVideoProcessQueueName: cfg.GetString("sqs_video_process_queue_name"),
			SQSVideoStatusQueueName:  cfg.GetString("sqs_video_status_queue_name"),
		}
	})

	return config
}

func initConfig() (viper.Viper, error) {
	cfg := viper.New()
	var err error
	initDefaults(cfg)
	// workaround because viper does not resolve envs when unmarshalling
	for _, key := range cfg.AllKeys() {
		val := cfg.Get(key)
		cfg.Set(key, val)
	}
	return *cfg, err
}

func initDefaults(config *viper.Viper) {
	config.SetDefault("server.host", "0.0.0.0:8000")
	config.SetDefault("sqs_video_process_queue_name", "video-process-queue")
	config.SetDefault("sqs_video_status_queue_name", "video-status-queue")
}
