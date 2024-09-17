package main

import (
	// "log"
	"os"
	"sync"

	"authen-service/appConfig/config"
	"authen-service/handlers"
	"authen-service/infrastructure"
)

func main() {
	env := ".env"

	// Check if there are command-line arguments
	if len(os.Args) > 1 {
		// If there is at least one argument, use it as the environment, name of file must be euqal to argument
		// exmapple: local env is local.env, the command must be `go run main local`
		env = os.Args[1]
	}
	config.SetEnvironment(env)

	cfg, _ := config.LoadConfig()
	db_authen, db_data, _ := infrastructure.Connect(cfg)
	redisCache := infrastructure.NewRedisClient(cfg)

	var wg sync.WaitGroup
	s, _ := handlers.NewServer(cfg, db_authen, db_data, redisCache)
	s.Start(&wg)
	wg.Wait()
}
