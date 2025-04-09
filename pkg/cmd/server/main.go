package main

import (
	server "github.com/MenD32/allpaca/pkg/server"
	config "github.com/MenD32/allpaca/pkg/server/config"
)

func main() {
	c := config.NewRecommendedConfig()
	s := server.NewServer(c)

	s.Start()
}
