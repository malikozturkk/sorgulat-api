package db

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var UniversityCollection *mongo.Collection

func ConnectMongoDB() {
    uri := "mongodb+srv://malikozturkk:CJwz4p9hSXkCEtVq@sorgulat-universities.engodae.mongodb.net/?retryWrites=true&w=majority&appName=sorgulat-universities"
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(uri)

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("MongoDB bağlantı hatası:", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("MongoDB erişim hatası:", err)
    }

    MongoClient = client
    UniversityCollection = client.Database("sorgulat").Collection("universities")

    log.Println("MongoDB bağlantısı başarılı.")
}
