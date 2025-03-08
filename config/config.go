package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Overload()
	if strings.Contains(string(err.Error()), "no such file or directory") {
		log.Printf("error loading .env file. Continue without it, gettings envs from environment...")
	} else {
		log.Fatalf("fail to read configs: %v", err)
		return
	}

}
