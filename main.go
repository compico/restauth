package main

import (
	"fmt"

	"github.com/compico/restauth/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	cfg      *db.MongoConfig
	database *db.DB
)

func init() {
	var err error
	fmt.Printf("[INFO] %v\n", "Creating MongoDB configuration")
	cfg, err = db.NewConfig()
	if err != nil {
		panic(err)
	}
	database = db.NewClient()
	fmt.Printf("[INFO] %v\n", "Creating new client for MongoDB")
	err = database.InitClient(cfg.URI)
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx, err := database.Connect()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		panic("")
	}
	defer func() {
		if err := database.Client.Disconnect(ctx); err != nil {
			fmt.Printf("[ERROR] %v\n", err.Error())
		}
	}()
	err = database.Ping()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		panic("")
	}
	coll := database.GetCollection("restauth", "testtoken")
	data := bson.M{
		"GUID":  "test1",
		"value": "test2",
	}
	res, err := database.InsertTransaction(coll, data)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		panic("")
	}
	fmt.Println(res.InsertedID)
}
