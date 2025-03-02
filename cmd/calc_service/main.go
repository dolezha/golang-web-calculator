package main

import (
	"calculator/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)
	http.HandleFunc("/api/v1/expressions", handlers.ExpressionsListHandler)
	http.HandleFunc("/api/v1/expressions/", handlers.ExpressionHandler)
	http.HandleFunc("/internal/task/", handlers.TaskHandler)
	http.HandleFunc("/internal/task", handlers.TaskHandler)

	log.Println("Оркестратор запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
