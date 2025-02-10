package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ApiAddress           string
	AccountAddress       string
	EmailAddress         string
	DSN                  string
	JWTSecret            string
	KafkaAddr            string
	EmailLogin           string
	EmailPass            string
	RoomAddress          string
	MessagesAddress      string
	WebsocketAddress     string
	NotificationsAddress string
}

func LoadConfig() *Config {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		panic(err)
	}
	return &Config{
		ApiAddress:           os.Getenv("ApiAddress"),
		AccountAddress:       os.Getenv("AccountAddress"),
		RoomAddress:          os.Getenv("RoomAddress"),
		MessagesAddress:      os.Getenv("MessagesAddress"),
		WebsocketAddress:     os.Getenv("WebsocketAddress"),
		NotificationsAddress: os.Getenv("NotificationsAddress"),
		EmailAddress:         os.Getenv("EmailAddress"),

		DSN:        os.Getenv("DSN"),
		JWTSecret:  os.Getenv("JWTSecret"),
		KafkaAddr:  os.Getenv("KafkaAddr"),
		EmailLogin: os.Getenv("EmailLogin"),
		EmailPass:  os.Getenv("EmailPass"),
	}
}
