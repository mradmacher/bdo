package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mradmacher/bdo/internal"
	"html/template"
	"os"
)

func loadData(filePath string) ([]bdo.Installation, error) {
	var installations []bdo.Installation
	jsonBlob, err := os.ReadFile(filePath)

	err = json.Unmarshal(jsonBlob, &installations)
	if err != nil {
		return nil, err
	}
	return installations, nil
}

func playWithDb() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	db := bdo.DbClient{}
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Disconnect(); err != nil {
			panic(err)
		}
	}()

	repo := db.NewInstallationRepo()

	result, err := repo.Find()
	if err != nil {
		panic(err)
	}
	fmt.Println("%v", *result)

	results, err := repo.Search(map[string]string{})
	if err != nil {
		panic(err)
	}
	for _, installation := range results {
		fmt.Println("%v", installation)
	}
}

type User struct {
	Name string
}

func playWithTemplates() {
	t, err := template.ParseFiles("cmd/playground/hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "John Smith",
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}

func main() {
	playWithTemplates()
}
