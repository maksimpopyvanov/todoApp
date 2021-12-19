package main

import (
	"log"
	todo "todo"
	handler "todo/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(todo.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Printf("error ocured while running http.server: %s", err.Error())
	}
}
