package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type qotdResp struct {
	Quote struct {
		Body string `json:"body"`
	} `json:"quote"`
}

func (s *Service) fetchQuote() (string, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	req, _ := http.NewRequest(http.MethodGet, "https://favqs.com/api/qotd", nil)

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Errorf("qotd request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("qotd non-200: %d", resp.StatusCode)
		return "", errors.New("qotd bad status")
	}

	body, _ := io.ReadAll(resp.Body)
	var data qotdResp
	if err := json.Unmarshal(body, &data); err != nil {
		s.logger.Errorf("qotd unmarshal: %v", err)
		return "", err
	}

	if data.Quote.Body == "" {
		return "", errors.New("qotd empty")
	}
	return data.Quote.Body, nil
}
