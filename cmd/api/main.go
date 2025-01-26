package main

import (
	"fmt"

	"github.com/8soat-grupo35/fastfood-order-production/internal/api/server"
	"github.com/8soat-grupo35/fastfood-order-production/internal/external"
)

func main() {
	fmt.Println("Iniciado o servidor Rest com GO")
	cfg := external.GetConfig()
	server.Start(cfg)
}
