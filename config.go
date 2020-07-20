package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
)

type Config struct {
	MongoDB MongoDB `json:"mongodb"`
}
type MongoDB struct {
	Hostname []string `json:"hostname"`
	DBname   string   `json:"dbname"`
	User     string   `json:"user"`
	Password string   `json:"password"`
}

/*
 JSON config "config.json" should look like this:
{
    "mongodb": {
        "hostname": [
            "hostname:27017",
            "hostname:27018",
            "hostname:27019"
        ],
        "dbname": "restauth",
        "user": "admin",
        "password": "123"
    }
}
*/

/*
Возвращает URI взятый из конфигов.
Работает только с серверами в режиме Replica Set
mongodb+srv://<username>:<password>@<hostname>/dbname?w=majority
*/
func (mongodb *MongoDB) getUri() (string, error) {
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
		return "", errors.New("Config error! Field 'HostName' is nil")
	}
	uri := url.URL{
		Scheme:   "mongodb",
		User:     user,
		Host:     mongodb.getHostname(),
		Path:     mongodb.DBname,
		RawQuery: "replicaSet=myrepls",
	}
	return uri.String(), nil
}

/*
Из массива адресов делает 1 большой адрес
соеденяя через запяту
*/
func (mongodb *MongoDB) getHostname() string {
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
Читает файл config.json и возвращает структуру
*/
func readConfig() (Config, error) {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return Config{}, err
	}
	config := new(Config)
	err = json.Unmarshal(file, config)
	if err != nil {
		return *config, err
	}
	return *config, nil
}
