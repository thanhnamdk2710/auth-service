package main

import (
	router "github.com/thanhnamdk2710/auth-service/internal/presentation/http"
)

func main() {
	server := router.NewRouter()

	server.Run(":8000")
}
