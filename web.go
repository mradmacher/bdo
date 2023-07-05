package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mradmacher/bdo/internal"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

  app := bdo.NewApp()
  app.MountHandlers()
  defer app.Stop()

	fmt.Println("Starting the server on :3000...")
  app.Start()
}
