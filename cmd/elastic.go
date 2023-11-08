package main

import (
	"elastic_web_service/internal/app"
	"log"
)

func main() {
	provider, err := app.NewProvider("./config/creds.yaml")
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}
	if provider.Repo().Init() != nil {
		log.Fatal("error in repo init")
	}
	log.Fatal(provider.Service().Run())
}
