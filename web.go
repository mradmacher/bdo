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
	renderer, err := bdo.NewRenderer("views")
	if err != nil {
		panic(err)
	}
	app, err := bdo.NewApp(*renderer)
	if err != nil {
		panic(err)
	}
	app.MountHandlers()
	defer app.Stop()

	fmt.Println("Starting the server on :3000...")
	app.Start()
}
