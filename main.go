package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    uri := os.Getenv("MONGODB_URI")
    fmt.Println("%s", uri)
    if uri == "" {
        log.Fatal("Set 'MONGODB_URI'")
    }
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        panic(err)
    }

    defer func() {
        if err := client.Disconnect(context.TODO()); err != nil {
            panic(err)
        }
    }()

    coll := client.Database("mbdo").Collection("installations")
    //_, err = coll.InsertOne(
    //    context.TODO(),
    //    bson.D{{"name", "Test"}})


    var result bson.M
    err = coll.FindOne(context.TODO(), bson.D{{"name", "Test"}}).Decode(&result)
    if err == mongo.ErrNoDocuments {
        fmt.Printf("No document\n")
        return
    }
    if err != nil {
        panic(err)
    }
    fmt.Println("%v", result)


    cursor, err := coll.Find(context.TODO(), bson.D{{"name", "Test"}})

    for cursor.Next(context.TODO()) {
        if err := cursor.Decode(&result); err != nil {
            panic(err)
        }
        fmt.Println("%v", result)

    }
    if err := cursor.Err(); err != nil {
        panic(err)
    }
}
