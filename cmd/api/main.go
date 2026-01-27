package main

import (
	"log"

	"github.com/thanhnamdk2710/auth-service/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to create application: " + err.Error())
	}

	if err := application.Run(); err != nil {
		log.Fatal("Application error: " + err.Error())
	}
}
