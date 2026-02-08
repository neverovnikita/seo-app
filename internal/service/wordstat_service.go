package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WordstatService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}
type WordstatResponse struct {
	Results []WordstatResult
}

type WordstatResult struct {
	RequestPhrase string
	TotalCount    *int64
	TopRequests   []PhraseCount
	Associations  []PhraseCount
	Error         *string
}

type PhraseCount struct {
	Phrase string
	Count  int64
}

func (wr *WordstatResponse) UnmarshalJSON(data []byte) error {
	var results []WordstatResult
	if err := json.Unmarshal(data, &results); err == nil {
		wr.Results = results
		return nil
	}

	var singleResult WordstatResult
	if err := json.Unmarshal(data, &singleResult); err == nil {
		wr.Results = []WordstatResult{singleResult}
		return nil
	}

	return json.Unmarshal(data, &wr.Results)
}

func NewWordstatService(baseURL, apiKey string, httpClient *http.Client) *WordstatService {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 120 * time.Second,
		}
	}
	return &WordstatService{baseURL, apiKey, httpClient}
}

type TopRequestsRequest struct {
	Phrase  string     `json:"phrase,omitempty"`  // одна фраза
	Phrases []string   `json:"phrases,omitempty"` // несколько фраз
	Regions []int      `json:"regions,omitempty"` // регионы (213 - Москва, 2 - СПб и т.д.)
	Devices []string   `json:"devices,omitempty"` // ["phone", "desktop", "tablet"]
	Date    *DateRange `json:"date,omitempty"`    // период
}
type DateRange struct {
	From string `json:"from"` // формат YYYY-MM-DD
	To   string `json:"to"`   // формат YYYY-MM-DD
}

func (s *WordstatService) RequestWordstat(req *TopRequestsRequest) (*WordstatResponse, error) {
	url := s.baseURL + "/topRequests"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json;charset=utf-8")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var wordstatResp WordstatResponse
	if err = json.Unmarshal(body, &wordstatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &wordstatResp, nil
}

func (s *WordstatService) Test() {

}
