package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mradmacher/bdo/internal"
	"os"
)

func loadData(filePath string, installations *[]bdo.Installation) error {
	jsonBlob, err := os.ReadFile(filePath)

	err = json.Unmarshal(jsonBlob, &installations)
	if err != nil {
		return err
	}
	return nil
}
func saveData(r *bdo.Repository, installations []bdo.Installation) error {
	for _, installation := range installations {
		id, err := installation.Add(r)
		fmt.Printf("%v, %v\n", id, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func try(err error) {
	if err != nil {
		panic(err)
	}
}

func seedDb() {
	db := &bdo.Repository{}
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Disconnect(); err != nil {
			panic(err)
		}
	}()

	var installations []bdo.Installation
	try(db.Purge())
	try(loadData("db/seed/installations.json", &installations))
	for _, installation := range installations {
		fmt.Printf("%v\n", installation)
	}
	try(saveData(db, installations))

}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	seedDb()
}
