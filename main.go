package main

import (
	"net/http"
	"time"

	"github.com/compico/restauth/internal/apiserver"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/api/v1/auth/:guid", apiserver.AddNewUser)
	router.GET("/api/v1/test/:guid", apiserver.Test)
	service := apiserver.ApiServer{
		Server: &http.Server{
			Addr:         ":8080",
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	err := service.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
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
