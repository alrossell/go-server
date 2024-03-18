package global

import (
    "fmt"
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const DataBase = "mySeriesDatabase"
const Collection = "mySeries"

var client *mongo.Client

func CreateClient() {
    fmt.Println("Creating client")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, _ = mongo.Connect(ctx, clientOptions)
}

func GetClient() *mongo.Client {
    if client == nil {
        CreateClient()
    }
    return client
}

func CloseClient() {
    fmt.Println("Closing client")
    if client != nil {
        client.Disconnect(context.Background())
    }
}

func GetDataBase() string {
    return DataBase
}

func GetCollection() string {
    return Collection
}
