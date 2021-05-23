package mongodb

import (
	"conkeys/config"
	"conkeys/storage"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct{}

type KeyStoreValue struct {
	Key   string
	Value storage.Value
}

var conn *mongo.Client
var db *mongo.Database

func getCtx(timeout time.Duration) (context.Context, context.CancelFunc) {
	t := 10 * time.Second
	if timeout > 0 {
		t = timeout
	}
	return context.WithTimeout(context.Background(), t)
}

func (m MongoStorage) Init() {
	cfg := config.GetConfig()
	connectionUri := "mongodb://127.0.0.1"
	if cfg.Mongo.ConnectionUri != "" {
		connectionUri = cfg.Mongo.ConnectionUri
	}
	fmt.Println(connectionUri)
	var err error
	conn, err = mongo.NewClient(options.Client().ApplyURI(connectionUri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := getCtx(0)
	defer cancel()

	err = conn.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db = conn.Database("keystorage")
}

func (m MongoStorage) Get(path string) (storage.Value, error) {

	col := db.Collection("keys")
	ctx, cancel := getCtx(0)
	defer cancel()
	filter := bson.D{
		primitive.E{
			Key:   "key",
			Value: path,
		},
	}
	var result KeyStoreValue
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return storage.Value{}, err
	}
	v := storage.Value{
		T: result.Value.T,
		V: result.Value.V,
	}
	return v, nil
}

func (m MongoStorage) GetKeys(pathSearch string) (map[string]storage.Value, error) {
	result := make(map[string]storage.Value)

	normalizedPathSearch := strings.TrimPrefix(pathSearch, "/")

	col := db.Collection("keys")
	ctx, cancel := getCtx(0)
	defer cancel()
	filter := bson.D{
		primitive.E{
			Key:   "key",
			Value: primitive.Regex{Pattern: fmt.Sprintf("^%s", normalizedPathSearch)},
		},
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", err)
	}

	for cursor.Next(ctx) {
		var element KeyStoreValue
		err := cursor.Decode(&element)
		if err != nil {
			fmt.Fprintf(gin.DefaultWriter, "%s", err)
		}
		result[element.Key] = element.Value
	}

	return result, nil
}

func (m MongoStorage) GetAllKeys() map[string]storage.Value {
	result := make(map[string]storage.Value)

	col := db.Collection("keys")
	ctx, cancel := getCtx(0)
	defer cancel()

	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s", err)
	}

	for cursor.Next(ctx) {
		var element KeyStoreValue
		err := cursor.Decode(&element)
		if err != nil {
			fmt.Fprintf(gin.DefaultWriter, "%s", err)
		}
		result[element.Key] = element.Value
	}

	return result
}

func (m MongoStorage) Put(path string, val storage.Value) {
	col := db.Collection("keys")
	ctx, cancel := getCtx(0)
	u := options.Replace()
	u.SetUpsert(true)
	defer cancel()

	normalizedPath := strings.TrimPrefix(path, "/")

	filter := bson.D{
		primitive.E{
			Key:   "key",
			Value: normalizedPath,
		},
	}

	v := KeyStoreValue{
		Key:   normalizedPath,
		Value: val,
	}

	col.ReplaceOne(ctx, filter, v, u)

}

func (m MongoStorage) Delete(path string) {
	col := db.Collection("keys")
	ctx, cancel := getCtx(0)
	defer cancel()
	col.DeleteOne(ctx, bson.M{"key": strings.TrimPrefix(path, "/")})
}
