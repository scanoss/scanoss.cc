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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type QueryParams map[string]string

type ScanossApiServiceHttpImpl struct {
	ctx     context.Context
	apiKey  string
	baseURL string
	client  HTTPClient
	timeout time.Duration
}

func NewScanossApiServiceHttpImpl() (ScanossApiService, error) {
	cfg := config.GetInstance()
	baseURL := cfg.GetApiUrl()
	apiKey := cfg.GetApiToken()

	if baseURL == "" {
		baseURL = config.SCANOSS_PREMIUM_API_URL
	}

	service := &ScanossApiServiceHttpImpl{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}

	return service, nil
}

func (s *ScanossApiServiceHttpImpl) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ScanossApiServiceHttpImpl) buildURL(endpoint string, params QueryParams) (string, error) {
	u, err := url.Parse(s.baseURL + endpoint)
	if err != nil {
		return "", err
	}

	if len(params) > 0 {
		q := u.Query()
		for key, value := range params {
			q.Add(key, value)
		}

		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

func (s *ScanossApiServiceHttpImpl) GetWithParams(endpoint string, params QueryParams) (*http.Response, error) {
	fullURL, err := s.buildURL(endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(s.ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
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

	params := QueryParams{
		"search":    request.Search,
		"vendor":    request.Vendor,
		"component": request.Component,
		"limit":     fmt.Sprintf("%d", request.Limit),
	}

	resp, err := s.GetWithParams("/api/v2/components/search", params)
	if err != nil {
		log.Error().Err(err).Msg("Error calling SCANOSS component search API")
		return entities.ComponentSearchResponse{}, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error().Int("statusCode", resp.StatusCode).Str("body", string(body)).Msg("API returned non-200 status")
		return entities.ComponentSearchResponse{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse entities.ComponentSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return entities.ComponentSearchResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Debug().
		Int("componentCount", len(apiResponse.Components)).
		Msg("Successfully retrieved component search results")

	return apiResponse, nil
}

// GetLicenseByPurl fetches licenses for a given purl/component combination using the SCANOSS HTTP API
func (s *ScanossApiServiceHttpImpl) GetLicensesByPurl(request entities.ComponentRequest) (entities.GetLicensesByPurlResponse, error) {
	if err := utils.GetValidator().Struct(request); err != nil {
		log.Error().Err(err).Msg("Invalid get license request")
		return entities.GetLicensesByPurlResponse{}, fmt.Errorf("invalid get license request: %w", err)
	}

	log.Debug().Str("purl", request.Purl).Msg("Fetching license for purl via SCANOSS API")

	if s.apiKey == "" {
		log.Error().Msg("SCANOSS API key not configured")
		return entities.GetLicensesByPurlResponse{}, fmt.Errorf("SCANOSS API key not configured")
	}

	params := QueryParams{
		"purl":        request.Purl,
		"requirement": request.Requirement,
	}

	resp, err := s.GetWithParams("/api/v2/licenses/component", params)
	if err != nil {
		log.Error().Err(err).Msg("Error calling SCANOSS component search API")
		return entities.GetLicensesByPurlResponse{}, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error().Int("statusCode", resp.StatusCode).Str("body", string(body)).Msg("API returned non-200 status")
		return entities.GetLicensesByPurlResponse{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse entities.GetLicensesByPurlResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return entities.GetLicensesByPurlResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return apiResponse, nil
}
