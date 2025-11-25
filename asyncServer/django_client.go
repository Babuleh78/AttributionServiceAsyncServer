package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type DjangoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDjangoClient(baseURL string) *DjangoClient {
	return &DjangoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetComposerAnalysis получает данные связи м-м по ID
func (c *DjangoClient) GetComposerAnalysis(composerAnalysisID int64, analysisID int64, composerID int64) (*ComposerAnalysis, error) {
	url := fmt.Sprintf("%s/api/attributionAnalyses/%d/composer/%d/", c.baseURL, analysisID, composerID)

	log.Printf("🔄 [GET] Fetching composer analysis from: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		log.Printf("❌ [GET] Failed to get composer analysis %d: %v", composerAnalysisID, err)
		return nil, fmt.Errorf("failed to get composer analysis: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("📥 [GET] Response status for composer analysis %d: %d", composerAnalysisID, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ [GET] Unexpected status code %d for composer analysis %d. Body: %s", resp.StatusCode, composerAnalysisID, string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var analysis ComposerAnalysis
	if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
		log.Printf("❌ [GET] Failed to decode composer analysis %d: %v", composerAnalysisID, err)
		return nil, fmt.Errorf("failed to decode composer analysis: %w", err)
	}

	log.Printf("✅ [GET] Successfully fetched composer analysis %d", composerAnalysisID)
	return &analysis, nil
}

// GetComposer получает данные композитора по ID
func (c *DjangoClient) GetComposer(composerID int64) (*Composer, error) {
	url := fmt.Sprintf("%s/api/composers/%d/", c.baseURL, composerID)

	log.Printf("🔄 [GET] Fetching composer from: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		log.Printf("❌ [GET] Failed to get composer %d: %v", composerID, err)
		return nil, fmt.Errorf("failed to get composer: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("📥 [GET] Response status for composer %d: %d", composerID, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ [GET] Unexpected status code %d for composer %d. Body: %s", resp.StatusCode, composerID, string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var composer Composer
	if err := json.NewDecoder(resp.Body).Decode(&composer); err != nil {
		log.Printf("❌ [GET] Failed to decode composer %d: %v", composerID, err)
		return nil, fmt.Errorf("failed to decode composer: %w", err)
	}

	log.Printf("✅ [GET] Successfully fetched composer %d: %s", composerID, composer.Name)
	return &composer, nil
}

// UpdateComposerAnalysis обновляет значение potential_coincidence
func (c *DjangoClient) UpdateComposerAnalysis(composerAnalysisID int64, analysisID int64, composerID int64, coincidence float64) error {
	url := fmt.Sprintf("%s/api/attributionAnalyses/%d/composer/%d/", c.baseURL, analysisID, composerID)

	log.Printf("🔄 [PUT] Updating composer analysis at: %s", url)

	updateReq := ComposerAnalysisUpdate{
		PotentialCoincidence: coincidence,
	}

	jsonData, err := json.Marshal(updateReq)
	if err != nil {
		log.Printf("❌ [PUT] Failed to marshal update request for %d: %v", composerAnalysisID, err)
		return fmt.Errorf("failed to marshal update request: %w", err)
	}

	log.Printf("📤 [PUT] Sending update data: %s", string(jsonData))

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ [PUT] Failed to create request for %d: %v", composerAnalysisID, err)
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("❌ [PUT] Failed to send update request for %d: %v", composerAnalysisID, err)
		return fmt.Errorf("failed to send update request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("📥 [PUT] Response status for update %d: %d", composerAnalysisID, resp.StatusCode)
	log.Printf("📥 [PUT] Response body for update %d: %s", composerAnalysisID, string(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ [PUT] Unexpected status code %d for update %d. Body: %s", resp.StatusCode, composerAnalysisID, string(body))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("✅ [PUT] Successfully updated composer analysis %d with coincidence: %.2f", composerAnalysisID, coincidence)
	return nil
}
