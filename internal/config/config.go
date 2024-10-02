package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env string
	Storage
	HttpServer
}

type Storage struct {
	Host     string
	Port     uint32
	User     string
	Password string
	DbName   string
}

type HttpServer struct {
	Address string
}

func ReadConfig(configPath string) Config {
	if configPath == "" {
		log.Fatalln("Config path is not set")
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Fatal(".env file load failed")
	}

	var cfg Config

	cfg.Env = os.Getenv("ENV")
	cfg.Storage.Host = os.Getenv("DB_HOST")

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("invalid .env file ", err.Error())
	}

	cfg.Storage.Port = uint32(port)
	cfg.Storage.User = os.Getenv("DB_USER")
	cfg.Storage.Password = os.Getenv("DB_PASSWORD")
	cfg.Storage.DbName = os.Getenv("DB_NAME")

	cfg.HttpServer.Address = os.Getenv("ADDR")

	return cfg
}
