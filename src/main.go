package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/prosperitybot/worker/app"
)

func main() {
	fmt.Println("Loading environment variables")
	_ = godotenv.Load()
	fmt.Println("Configuring Logging system...")
	app.Start()
}
