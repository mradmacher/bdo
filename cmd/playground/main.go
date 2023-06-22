package main

import (
    "fmt"
    "os"
    "encoding/json"
    "github.com/joho/godotenv"
    "github.com/mradmacher/mbdo/internal"
)

func loadData(filePath string) ([]bdo.Installation, error) {
    var installations []bdo.Installation
    jsonBlob, err := os.ReadFile(filePath)

    err = json.Unmarshal(jsonBlob, &installations)
    if err != nil { return nil, err }
    return installations, nil
}

func playWithDb() {
    db := bdo.DbClient{}
    err := db.Connect()
    if err != nil { panic(err) }

    defer func() {
        if err := db.Disconnect(); err != nil {
            panic(err)
        }
    }()

    repo := db.NewInstallationRepo()
    err = repo.Purge()
    if err != nil { panic(err) }

    var installations []bdo.Installation
    installations, err = loadData("db_seed.json")
    if err != nil { panic(err) }

    for _, installation := range installations {
        err = repo.Add(&installation)
        if err != nil { panic(err) }
    }

    result, err := repo.Find()
    if err != nil { panic(err) }
    fmt.Println("%v", *result)

    results, err := repo.Search(map[string]string{})
    if err != nil { panic(err) }
    for _, installation := range results {
        fmt.Println("%v", installation)
    }
}

func main() {
    if err := godotenv.Load(); err != nil {
        panic("No .env file found")
    }

    playWithDb()
}
