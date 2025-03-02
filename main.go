package main

import (
	"calculator/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)

	fmt.Println("Сервер запущен на 8080 порту")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
