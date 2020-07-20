package main

import "fmt"

var config Config

func init() {
	var err error
	config, err = readConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Printf("%#v", config)
}
