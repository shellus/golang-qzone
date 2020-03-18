package main

import (
	"./qzone"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("Error loading .env file: %s", err))
	}

	qzone.Login(
		os.Getenv("COOKIESTR"),
		os.Getenv("G_TK"),
		os.Getenv("QQ"),
		os.Getenv("TOPICID"),
		os.Getenv("DOWNLOADDIR"))
	qzone.Run()
}
