package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struct yang digunakan untuk menyimpan informasi database
type DBInfo struct {
	DBString string
	DBName   string
}

// MongoConnect adalah fungsi untuk menghubungkan ke database MongoDB
func MongoConnect(mconn DBInfo) (db *mongo.Database, err error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mconn.DBString))
	if err != nil {
		return nil, err
	}
	return client.Database(mconn.DBName), nil
}

var MongoString string = GetEnv("MONGOSTRING")

var mongoinfo = DBInfo{
	DBString: MongoString,
	DBName:   "db_barokah",
}

var Mongoconn, ErrorMongoconn = MongoConnect(mongoinfo)
