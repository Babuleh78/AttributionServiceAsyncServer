package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type AnalysisHandler struct {
	djangoClient *DjangoClient
	calculator   *CoincidenceCalculator
	mu           sync.Mutex
	processing   map[int64]bool
}

func NewAnalysisHandler(djangoClient *DjangoClient, calculator *CoincidenceCalculator) *AnalysisHandler {
	return &AnalysisHandler{
		djangoClient: djangoClient,
		calculator:   calculator,
		processing:   make(map[int64]bool),
	}
}

// CalculateCoincidenceHandler асинхронный обработчик для расчета совпадения
func (h *AnalysisHandler) CalculateCoincidenceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ComposerAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем секретный ключ (хардкод)
	expectedSecretKey := "music_analysis_secret_2024"
	if req.SecretKey != expectedSecretKey {
		http.Error(w, "Invalid secret key", http.StatusUnauthorized)
		return
	}

	// Проверяем, не обрабатывается ли уже этот анализ
	h.mu.Lock()
	if h.processing[req.ComposerAnalysisID] {
		h.mu.Unlock()
		http.Error(w, "Analysis already being processed", http.StatusConflict)
		return
	}
	h.processing[req.ComposerAnalysisID] = true
	h.mu.Unlock()

	// Асинхронная обработка
	go h.processComposerAnalysis(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"message": "Coincidence calculation started",
	})
}

// CalculateCoincidenceSyncHandler СИНХРОННЫЙ обработчик для расчета совпадения
func (h *AnalysisHandler) CalculateCoincidenceSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ComposerAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем секретный ключ (хардкод)
	expectedSecretKey := "music_analysis_secret_2024"
	if req.SecretKey != expectedSecretKey {
		http.Error(w, "Invalid secret key", http.StatusUnauthorized)
		return
	}

	log.Printf("🔄 [SYNC] Processing composer analysis ID: %d", req.ComposerAnalysisID)

	// СИНХРОННО обрабатываем запрос
	result, err := h.processComposerAnalysisSync(req)
	if err != nil {
		log.Printf("❌ [SYNC] Error processing composer analysis %d: %v", req.ComposerAnalysisID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// processComposerAnalysis асинхронно обрабатывает расчет совпадения
func (h *AnalysisHandler) processComposerAnalysis(req ComposerAnalysisRequest) {
	defer func() {
		h.mu.Lock()
		delete(h.processing, req.ComposerAnalysisID)
		h.mu.Unlock()
	}()

	// Получаем данные связи м-м
	composerAnalysis, err := h.djangoClient.GetComposerAnalysis(req.ComposerAnalysisID, req.AnalysisID, req.ComposerID)
	if err != nil {
		log.Printf("❌ [PROCESS] Failed to get composer analysis %d: %v", req.ComposerAnalysisID, err)
		return
	}

	// Получаем данные композитора
	composer, err := h.djangoClient.GetComposer(req.ComposerID)
	if err != nil {
		log.Printf("❌ [PROCESS] Failed to get composer %d: %v", req.ComposerID, err)
		return
	}

	// Рассчитываем совпадение
	coincidence := h.calculator.CalculatePotentialCoincidence(composer, composerAnalysis)

	// Отправляем результат обратно в Django
	h.sendResultToDjango(req.ComposerAnalysisID, coincidence)
}

// processComposerAnalysisSync СИНХРОННО обрабатывает расчет совпадения
func (h *AnalysisHandler) processComposerAnalysisSync(req ComposerAnalysisRequest) (*ComposerAnalysisSyncResponse, error) {
	// Получаем данные связи м-м
	composerAnalysis, err := h.djangoClient.GetComposerAnalysis(req.ComposerAnalysisID, req.AnalysisID, req.ComposerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get composer analysis: %w", err)
	}

	// Получаем данные композитора
	composer, err := h.djangoClient.GetComposer(req.ComposerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get composer: %w", err)
	}

	// Рассчитываем совпадение
	coincidence := h.calculator.CalculatePotentialCoincidence(composer, composerAnalysis)

	// Возвращаем результат
	result := &ComposerAnalysisSyncResponse{
		ComposerAnalysisID:   req.ComposerAnalysisID,
		PotentialCoincidence: coincidence,
		SecretKey:            "music_analysis_secret_2024",
		Status:               "completed",
	}

	log.Printf("✅ [SYNC] Successfully calculated coincidence for composer analysis %d: %.2f%%",
		req.ComposerAnalysisID, coincidence)

	return result, nil
}

// sendResultToDjango отправляет результат расчета в Django (для асинхронного режима)
func (h *AnalysisHandler) sendResultToDjango(composerAnalysisID int64, coincidence float64) {
	result := ComposerAnalysisResponse{
		ComposerAnalysisID:   composerAnalysisID,
		PotentialCoincidence: coincidence,
		SecretKey:            "music_analysis_secret_2024",
	}

	// URL для callback в Django
	callbackURL := h.djangoClient.baseURL + "/api/analysis-callback/"

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("❌ [CALLBACK] Failed to marshal result: %v", err)
		return
	}

	resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ [CALLBACK] Failed to send result to Django: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ [CALLBACK] Django returned status %d", resp.StatusCode)
		return
	}

	log.Printf("✅ [CALLBACK] Successfully sent result to Django for analysis %d", composerAnalysisID)
}
