package main

import (
	"github.com/Imm0bilize/gunshot-api-service/internal/app"
	"github.com/Imm0bilize/gunshot-api-service/internal/config"
	"log"
)

func main() {
	cfg, err := config.New(".env.public", ".env.private")
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
