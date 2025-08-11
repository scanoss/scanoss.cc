// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2025 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/fetch"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ScanossApiServiceHttpImpl struct {
	client   *http.Client
	timeout  time.Duration
	endpoint string
	apiKey   string
}

func NewScanossApiServiceHttpImpl() (ScanossApiService, error) {
	cfg := config.GetInstance()

	endpoint := cfg.GetApiUrl()
	if endpoint == "" {
		endpoint = config.SCANOSS_PREMIUM_API_URL
	}

	service := &ScanossApiServiceHttpImpl{
		timeout:  30 * time.Second,
		endpoint: endpoint,
		apiKey:   cfg.GetApiToken(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return service, nil
}

// SearchComponents searches for components using the SCANOSS HTTP API
func (s *ScanossApiServiceHttpImpl) SearchComponents(request entities.ComponentSearchRequest) (entities.ComponentSearchResponse, error) {
	if err := utils.GetValidator().Struct(request); err != nil {
		log.Error().Err(err).Msg("Invalid component search request")
		return entities.ComponentSearchResponse{}, fmt.Errorf("invalid search request: %w", err)
	}

	log.Debug().
		Str("search", request.Search).
		Str("vendor", request.Vendor).
		Str("component", request.Component).
		Int32("limit", request.Limit).
		Msg("Searching components via SCANOSS API")

	if s.apiKey == "" {
		log.Error().Msg("SCANOSS API key not configured")
		return entities.ComponentSearchResponse{}, fmt.Errorf("SCANOSS API key not configured")
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling request body")
		return entities.ComponentSearchResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := fmt.Sprintf("%s/api/v2/components/search", s.endpoint)

	options := fetch.Options{
		Method: http.MethodPost,
		Headers: map[string]string{
			"x-api-key":    s.apiKey,
			"Content-Type": "application/json",
		},
		Body: jsonBody,
	}

	resp, err := fetch.Fetch(apiURL, options)
	if err != nil {
		log.Error().Err(err).Msg("Error calling SCANOSS component search API")
		return entities.ComponentSearchResponse{}, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading API response")
		return entities.ComponentSearchResponse{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("statusCode", resp.StatusCode).
			Str("body", string(body)).
			Msg("API returned non-200 status")
		return entities.ComponentSearchResponse{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse entities.ComponentSearchResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Error parsing API response")
		return entities.ComponentSearchResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	components := make([]entities.SearchedComponent, len(apiResponse.Components))
	for i, comp := range apiResponse.Components {
		components[i] = entities.SearchedComponent{
			Component: comp.Component,
			Purl:      comp.Purl,
			URL:       comp.URL,
		}
	}

	response := entities.ComponentSearchResponse{
		Components: components,
		Status: entities.StatusResponse{
			Status:  apiResponse.Status.Status,
			Message: apiResponse.Status.Message,
		},
	}

	log.Debug().
		Int("componentCount", len(response.Components)).
		Msg("Successfully retrieved component search results")

	return response, nil
}

// Close closes the HTTP client (no-op for HTTP client)
