package tpl

var (
	MigrateMongoGo = `package mongo

import (
	"fmt"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	name := "add_index_to_users"
	version := "1642035425"

	migrate.Register(
		func(db *mongo.Database) error {
			fmt.Printf("[UP] %s - version %s\n", name, version)

			return nil
		},
		func(db *mongo.Database) error {
			fmt.Printf("[DOWN] %s - version %s\n", name, version)

			return nil
		},
	)
}
`

	MigrateMongoMainGo = `package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	mongoDriver "github.com/Drafteame/framework/mongo"
	"github.com/spf13/viper"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"

	_ "{{Namespace}}/migrations/mongo/migrate"
)

func init() {
	viper.SetDefault("STAGE", "local")
}

func main() {
	defer func() {
		if recov := recover(); recov != nil {
			fmt.Println("Error:", recov)
			fmt.Println(string(debug.Stack()))
		}
	}()

	viper.AutomaticEnv()
	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	migrate.SetDatabase(db)

	action := os.Args[1]

	fmt.Printf("Mongo migrate %s: \n\n", action)

	switch action {
	case "up":
		if err := migrate.Up(migrate.AllAvailable); err != nil {
			panic(err)
		}
	case "down":
		if err := migrate.Down(migrate.AllAvailable); err != nil {
			panic(err)
		}
	default:
		panic("unknown action")
	}
}

func connectDB() (*mongo.Database, error) {
	var config mongoDriver.Config

	if viper.GetString("STAGE") == "local" {
		viper.AddConfigPath(".")
		viper.SetConfigName(".engine")
		viper.SetConfigType("yaml")

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}

		config = mongoDriver.Config{URI: viper.GetString("MONGO_URI")}
	} else {
		config = mongoDriver.Config{
			ReadPreference:  viper.GetString("MONGO_READ_PREFERENCE"),
			UserName:        viper.GetString("MONGO_USERNAME"),
			Password:        viper.GetString("MONGO_PASSWORD"),
			ClusterEndpoint: viper.GetString("MONGO_CLUSTER_ENDPOINT"),
			CertPath:        viper.GetString("MONGO_CERT_PATH"),
			DBName:          viper.GetString("MONGO_DB_NAME"),
		}
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Config: %s\nDatabase: %s\n\n", string(configBytes), viper.GetString("MONGO_DB_NAME"))

	client, err := mongoDriver.NewWithConfig(config)
	if err != nil {
		return nil, err
	}

	db := client.Client().Database(viper.GetString("MONGO_DB_NAME"))

	return db, nil
}
`
)
