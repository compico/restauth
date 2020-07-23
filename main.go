package main

import (
	"fmt"

	"github.com/compico/restauth/internal/db"
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
	fmt.Printf("[DEBUG] Your URI: %v\n", cfg.URI)
	database = db.NewClient()
	fmt.Printf("[INFO] %v\n", "Creating new client for MongoDB")
	err = database.InitClient(cfg.URI)
	if err != nil {
		panic(err)
	}
}

func main() {

}

/*
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
	err = database.AddUserInDB("restauth", "c3b6de2f-756b-41ce-9072-9fd5210de008")
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		panic(err)
	}
	coll := database.GetCollection("restauth", "c3b6de2f-756b-41ce-9072-9fd5210de008")
	data := bson.M{
		"GUID":  "test1",
		"value": "test2",
	}
	restrans, err := database.InsertTransaction(coll, data)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		panic(err)
	}
	fmt.Println(restrans.InsertedID)
*/
