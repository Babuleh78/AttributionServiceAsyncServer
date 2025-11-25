package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	// Настройка логгера для вывода времени и даты
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	// Получаем значения из переменных окружения или используем по умолчанию
	djangoURL := getEnv("DJANGO_API_URL", "http://django:8000")
	port := getEnv("PORT", "8888")

	log.Printf("🚀 ============ Starting Music Analysis Service ============")
	log.Printf("⚙️  Configuration:")
	log.Printf("   Django API URL: %s", djangoURL)
	log.Printf("   Service Port: %s", port)

	// Инициализация компонентов
	log.Printf("🔧 Initializing components...")
	djangoClient := NewDjangoClient(djangoURL)
	calculator := NewCoincidenceCalculator()
	handler := NewAnalysisHandler(djangoClient, calculator)

	// Настройка маршрутов
	http.HandleFunc("/api/calculate-coincidence", handler.CalculateCoincidenceHandler)
	http.HandleFunc("/api/calculate-coincidence-sync", handler.CalculateCoincidenceSyncHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	log.Printf("🌐 Server starting on port %s", port)
	log.Printf("📡 Endpoints:")
	log.Printf("   POST /api/calculate-coincidence")
	log.Printf("   POST /api/calculate-coincidence-sync")
	log.Printf("   GET  /health")
	log.Printf("======================================================")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("💥 Failed to start server: %v", err)
		os.Exit(1)
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
