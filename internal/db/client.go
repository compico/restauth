package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Наследование клиента из драйвера для создания новых методов
type DB struct {
	Client *mongo.Client
}

func NewClient() *DB {
	db := new(DB)
	return db
}

//Новый клиент для MongoDB
func (db *DB) InitClient(uri string) error {
	var err error

	db.Client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	return nil
}

//Подключение серверу
func (db *DB) Connect() (context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := db.Client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

//Отправка пинг сигнала для проверки состояния сервера
func (db *DB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}

/*Получение коллекции из базы данных.
На вход название базы, название колекции.
На выходе коллекция из базы.*/
func (db *DB) GetCollection(dbName, collection string) *mongo.Collection {
	coll := db.Client.Database(dbName).Collection(collection)
	return coll
}

//Создание сессии для будущих транзакций
func (db *DB) StartSession() (mongo.Session, error) {
	session, err := db.Client.StartSession()
	if err != nil {
		return nil, err
	}
	return session, err
}

//Тестовая функция для проверки
func InsertOne(coll *mongo.Collection) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := coll.InsertOne(ctx, bson.M{
		"GUID":  "2ADEA345-7F7A-4313-87AE-F05E8B2DE678",
		"value": "Yzg0ZjE4YTItYzZjNy00ODUwLWJlMTUtOTNmOWNiYWVmM2Iz",
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
