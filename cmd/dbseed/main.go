package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mradmacher/bdo/internal/repo"
	"os"
)

func loadData(filePath string, installations *[]repo.Installation) error {
	jsonBlob, err := os.ReadFile(filePath)

	err = json.Unmarshal(jsonBlob, &installations)
	if err != nil {
		return err
	}
	return nil
}
func saveData(r *repo.Repository, installations []repo.Installation) (error) {
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
	db := &repo.Repository{}
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Disconnect(); err != nil {
			panic(err)
		}
	}()

	var installations []repo.Installation
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
