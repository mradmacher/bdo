package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mradmacher/bdo/internal"
	"bytes"
	"os/exec"
	"runtime"
	_ "embed"
)

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		fmt.Printf("Failed to open browser: %v", err)
	}
}

//go:embed .env
var envFile []byte

func main() {
	envMap, err := godotenv.Parse(bytes.NewReader(envFile))
	if err != nil {
		panic("No .env file found")
	}
	renderer, err := bdo.NewRenderer()
	if err != nil {
		panic(err)
	}
	dbUri := envMap["BDO_DB_URI"]
	mapsApiKey := envMap["GOOGLE_MAPS_API_KEY"]
	fmt.Printf("Loading data from %s\n", dbUri)

	app, err := bdo.NewApp(*renderer, dbUri, mapsApiKey)
	if err != nil {
		panic(err)
	}
	app.MountHandlers()
	defer app.Stop()

	fmt.Println("Starting the server on :3001...")
	go func() {
		app.Start()
	}()
	openBrowser("http://localhost:3001")

	select {}
}
