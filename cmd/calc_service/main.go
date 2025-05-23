package main

import (
	"calculator/handlers"
	"calculator/middleware"
	"calculator/services"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	dbPath := getEnv("DB_PATH", "./calculator.db")
	db, err := services.NewDatabaseService(dbPath)
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	defer db.Close()

	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	authService := services.NewAuthService(db, jwtSecret)
	expressionService := services.NewExpressionService(db)

	authHandler := handlers.NewAuthHandler(authService)
	calculateHandler := handlers.NewCalculateHandler(expressionService)
	expressionHandler := handlers.NewExpressionHandler(expressionService)
	taskHandler := handlers.NewTaskHandler(expressionService)

	authMiddleware := middleware.AuthMiddleware(authService)

	http.HandleFunc("/api/v1/register", authHandler.Register)
	http.HandleFunc("/api/v1/login", authHandler.Login)

	http.HandleFunc("/internal/task", taskHandler.GetTask)
	http.HandleFunc("/internal/task/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskHandler.GetTaskByID(w, r)
		} else if r.Method == http.MethodPost {
			taskHandler.SubmitTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/api/v1/calculate", authMiddleware(http.HandlerFunc(calculateHandler.Calculate)))
	http.Handle("/api/v1/expressions", authMiddleware(http.HandlerFunc(expressionHandler.GetExpressions)))
	http.Handle("/api/v1/expressions/", authMiddleware(http.HandlerFunc(expressionHandler.GetExpression)))

	port := getEnv("PORT", "8080")
	fmt.Printf("Server started on port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
