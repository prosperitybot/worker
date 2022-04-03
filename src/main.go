package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/prosperitybot/worker/app"
	"github.com/prosperitybot/worker/logging"
)

func main() {
	fmt.Println("Loading environment variables")
	_ = godotenv.Load()
	fmt.Println("Configuring Logging system...")
	_ = logging.Init()
	app.Start()
}
