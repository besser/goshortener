package database

import (
	"fmt"
	"os"

	"workspace/src/gopkg.in/mgo.v2"

	"github.com/besser/goshortener/config"
)

func ConnectarMongoDB() *mgo.Session {
	uri := config.Cfg.GetString("database.mongo_uri")

	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	return sess
}
