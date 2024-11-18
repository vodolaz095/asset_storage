package config

import "os"

// DSN is data source name used to dial PostgeSQL database
var DSN = "user=assets password=assets dbname=assets host=localhost port=5432 sslmode=disable"

// Address shows where application binds
var Address = "0.0.0.0"

// Port sets port number application binds onto
var Port = "3000"

func loadFromEnvironment(v *string, key string) {
	if fromEnv := os.Getenv(key); fromEnv != "" {
		*v = fromEnv
	}
}

func Load() {
	loadFromEnvironment(&Port, "PORT")
	loadFromEnvironment(&Address, "ADDR")
	loadFromEnvironment(&DSN, "DB_URL")
}
