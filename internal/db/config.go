package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

type MongoConfig struct {
	Hostname   []string `json:"hostname"`
	DBname     string   `json:"dbname"`
	User       string   `json:"user"`
	Password   string   `json:"password"`
	ReplicaSet string   `json:"replicaset"`
	URI        string   `json:"-"`
}

/*
 JSON конфиг "mongo.json"
 должен лежать в папке cfg
 И должен выглядеть так:
{
 "hostname": [
        "hostname:27017",
        "hostname:27018",
        "hostname:27019"
	],
	"replicaset": "namereplicaset"
    "dbname": "restauth",
    "user": "admin",
    "password": "123"
}

*/

func NewConfig() (*MongoConfig, error) {
	mongodb := new(MongoConfig)
	fmt.Printf("[INFO] %v\n", "Reading config.")
	err := mongodb.readConfig()
	if err != nil {
		return nil, err
	}
	fmt.Printf("[INFO] %v\n", "Config ready.")
	fmt.Printf("[INFO] %v\n", "Getting URI for connetion.")
	err = mongodb.getUri()
	if err != nil {
		return nil, err
	}
	fmt.Printf("[INFO] %v\n", "URI ready.")
	return mongodb, nil
}

/*
Возвращает URI взятый из конфигов.
Работает только с серверами в режиме Replica Set
mongodb+srv://<username>:<password>@<hostname>/dbname?w=majority
*/
func (mongodb *MongoConfig) getUri() error {
	//"mongodb://<username>:<password>@<cluster-address>/test?w=majority"
	var user *url.Userinfo
	if mongodb.User != "" {
		if mongodb.Password != "" {
			user = url.UserPassword(mongodb.User, mongodb.Password)
		}
		if mongodb.Password == "" {
			user = url.User(mongodb.User)
		}
	}
	if len(mongodb.Hostname) == 0 {
		return errors.New("Config error! Field 'HostName' is nil")
	}
	query := "authSource=admin&replicaSet=" +
		mongodb.ReplicaSet + "&readPreference=primary&ssl=true"
	x := url.URL{
		Scheme:   "mongodb+srv",
		User:     user,
		Host:     mongodb.getHostname(),
		Path:     "/admin",
		RawQuery: query,
	}
	mongodb.URI = x.String()
	return nil
}

/*
Из массива адресов делает 1 большой адрес
соеденяя через запяту
*/
func (mongodb *MongoConfig) getHostname() string {
	if len(mongodb.Hostname) == 1 {
		return mongodb.Hostname[0]
	}
	var hosts string
	for i := 0; i < len(mongodb.Hostname); i++ {
		hosts += mongodb.Hostname[i]
		if i == len(mongodb.Hostname)-1 {
			break
		}
		hosts += ","
	}
	return hosts
}

/*
Читает файл mongo.json и возвращает структуру
*/
func (config *MongoConfig) readConfig() error {
	file, err := ioutil.ReadFile("./cfg/mongo.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		return err
	}
	return nil
}
