package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	err := db.Client.Ping(ctx, options.Client().ReadPreference)
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

func (db *DB) AddUserInDB(dbName, guid string) error {
	session, err := db.Client.StartSession()
	if err != nil {
		return err
	}
	err = session.StartTransaction()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err = session.CommitTransaction(sc)
		if err != nil {
			return err
		}
		err := db.Client.Database(dbName).CreateCollection(sc, guid)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	session.EndSession(ctx)
	return nil
}

//Метод записи данных в транзакции
func (db *DB) InsertTransaction(coll *mongo.Collection, data bson.M) (*mongo.InsertOneResult, error) {
	session, err := db.Client.StartSession()
	if err != nil {
		return nil, err
	}
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	var res *mongo.InsertOneResult
	err = mongo.WithSession(context.Background(), session,
		func(sc mongo.SessionContext) error {
			err = session.CommitTransaction(sc)
			if err != nil {
				return err
			}
			res, err = coll.InsertOne(sc, data)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	session.EndSession(context.Background())
	return res, nil
}
