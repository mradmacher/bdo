package main

import (
    "fmt"
    "os"
    "encoding/json"
    "github.com/joho/godotenv"
    "github.com/mradmacher/mbdo/internal/repo"
)

func loadData(filePath string) ([]repo.Installation, error) {
    var installations []repo.Installation
    jsonBlob, err := os.ReadFile(filePath)

    err = json.Unmarshal(jsonBlob, &installations)
    if err != nil { return nil, err }
    return installations, nil
}

func playWithDb() {
    db := repo.DbClient{}
    err := db.Connect()
    if err != nil { panic(err) }

    defer func() {
        if err := db.Disconnect(); err != nil {
            panic(err)
        }
    }()

    installationsRepo := db.NewInstallationRepo()
    err = installationsRepo.Purge()
    if err != nil { panic(err) }

    var installations []repo.Installation
    installations, err = loadData("db_seed.json")
    if err != nil { panic(err) }

    for _, installation := range installations {
        err = installationsRepo.Add(&installation)
        if err != nil { panic(err) }
    }

    result, err := installationsRepo.Find()
    if err != nil { panic(err) }
    fmt.Println("%v", *result)

    results, err := installationsRepo.Search()
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
