package main

import (
    "context"
    "fmt"
    "errors"
    "os"
    "encoding/json"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func connectClient() (*mongo.Client, error) {
    if err := godotenv.Load(); err != nil {
        return nil, errors.New("No .env file found")
    }

    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        return nil, errors.New("Set 'MONGODB_URI'")
    }
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }

    return client, nil
}

type Decision struct {
    WasteCode string `bson:"waste_code"`
    ProcessCode string `bson:"process_code,omitempty"`
    Quantity int `bson:"quantity,omitempty"`
}

type Address struct {
    Line1 string
    Line2 string
    StateCode string
}

type Installation struct {
    Name string
    Nip string
    Regon string
    Address Address
    Decisions []Decision
}

type InstallationRepo struct {
    Collection *mongo.Collection
}

func (repo *InstallationRepo) purge() error {
    _, err := repo.Collection.DeleteMany(context.TODO(), bson.D{})
    return err
}

func (repo *InstallationRepo) add(installation *Installation) error {
    _, err := repo.Collection.InsertOne(context.TODO(), installation)
    if err != nil { return err }
    return nil
}

func (repo *InstallationRepo) search(installations *[]Installation) error {
    cursor, err := repo.Collection.Find(context.TODO(), bson.D{})
    if err != nil { return err }
    for cursor.Next(context.TODO()) {
        var result Installation
        if err := cursor.Decode(&result); err != nil {
            return err
        }
        *installations = append(*installations, result)
    }
    if err := cursor.Err(); err != nil {
        return err
    }
    return nil
}

func (repo *InstallationRepo) find(installation *Installation) error {
    err := repo.Collection.FindOne(context.TODO(), bson.D{}).Decode(installation)
    if err == mongo.ErrNoDocuments {
        installation = nil
        return nil
    } else if err != nil {
        return err
    }
    return nil
}

func loadData(filePath string, installations *[]Installation) error {
    jsonBlob, err := os.ReadFile(filePath)

    err = json.Unmarshal(jsonBlob, installations)
    if err != nil { return err }
    return nil
}

func main() {
    client, err := connectClient()
    if err != nil { panic(err) }

    defer func() {
        if err := client.Disconnect(context.TODO()); err != nil {
            panic(err)
        }
    }()

    repo := InstallationRepo{Collection: client.Database("mbdo").Collection("installations")}
    err = repo.purge()
    if err != nil { panic(err) }

    var installations []Installation
    err = loadData("db_seed.json", &installations)
    if err != nil { panic(err) }

    for _, installation := range installations {
        err = repo.add(&installation)
        if err != nil { panic(err) }
    }

    var result Installation
    err = repo.find(&result)
    if err != nil { panic(err) }
    fmt.Println("%v", result)

    var results []Installation
    err = repo.search(&results)
    if err != nil { panic(err) }
    for _, installation := range results {
        fmt.Println("%v", installation)
    }
}
