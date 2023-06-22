package bdo

import (
    "context"
    "errors"
    "os"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type DbClient struct {
    Client *mongo.Client
}

func (db *DbClient) Connect() error {
    var err error
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        return errors.New("Set 'MONGODB_URI'")
    }
    db.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        return err
    }

    return nil
}

func (db *DbClient) Disconnect() error {
    return db.Client.Disconnect(context.TODO())
}

func (db *DbClient) NewInstallationRepo() *InstallationRepo {
    return &InstallationRepo{Collection: db.Client.Database("mbdo").Collection("installations")}
}

type Params map[string]string

type Capability struct {
    WasteCode string `bson:"waste_code"`
    Dangerous bool `bson:"dangerous"`
    ProcessCode string `bson:"process_code,omitempty"`
    Quantity int `bson:"quantity,omitempty"`
}

type Address struct {
    Line1 string
    Line2 string
    StateCode string `bson:"state_code"`
}

type Installation struct {
    Name string
    Address Address
    Capabilities []Capability
}

type InstallationRepo struct {
    Collection *mongo.Collection
}

func (repo *InstallationRepo) Purge() error {
    _, err := repo.Collection.DeleteMany(context.TODO(), bson.D{})
    return err
}

func (repo *InstallationRepo) Add(installation *Installation) error {
    _, err := repo.Collection.InsertOne(context.TODO(), installation)
    if err != nil { return err }
    return nil
}

func (repo *InstallationRepo) Search(params Params) ([]Installation, error) {
    var installations []Installation

    query := bson.D{}
    for k, v := range params {
        switch k {
            case "process_code":
                query = append(query, bson.E{"capabilities.process_code", v})
            case "waste_code":
                query = append(query, bson.E{"capabilities.waste_code", v})
        }
    }
    cursor, err := repo.Collection.Find(context.TODO(), query)
    if err != nil { return nil, err }

    for cursor.Next(context.TODO()) {
        var result Installation
        if err := cursor.Decode(&result); err != nil {
            return nil, err
        }
        installations = append(installations, result)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    return installations, nil
}

func (repo *InstallationRepo) Find() (*Installation, error) {
    var installation Installation
    err := repo.Collection.FindOne(context.TODO(), bson.D{}).Decode(&installation)
    if err == mongo.ErrNoDocuments {
        return nil, nil
    } else if err != nil {
        return nil, err
    }
    return &installation, nil
}
