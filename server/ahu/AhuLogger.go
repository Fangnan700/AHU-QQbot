package ahu

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(info string) {
	logger.SetPrefix("[INFO]")
	logger.Println(info)
}

func LogError(err error) {
	logger.SetPrefix("[ERROR]")
	logger.Println(err)
}
