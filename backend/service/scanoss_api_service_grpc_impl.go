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
	"crypto/tls"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/papi/api/componentsv2"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/mappers"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type ScanossApiServiceGrpcImpl struct {
	client   componentsv2.ComponentsClient
	conn     *grpc.ClientConn
	mapper   mappers.ScanossApiMapper
	timeout  time.Duration
	endpoint string
	apiKey   string
}

func NewScanossApiServiceGrpcImpl(
	mapper mappers.ScanossApiMapper,
) ScanossApiService {
	cfg := config.GetInstance()
	service := &ScanossApiServiceGrpcImpl{
		mapper:   mapper,
		timeout:  30 * time.Second,
		endpoint: utils.RemoveProtocolFromUrl(cfg.GetApiUrl()),
		apiKey:   cfg.GetApiToken(),
	}

	if err := service.initializeGrpcClient(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize gRPC client")
	}

	return service
}

// SearchComponents searches for components using the SCANOSS gRPC API
func (s *ScanossApiServiceGrpcImpl) SearchComponents(request entities.ComponentSearchRequest) (entities.ComponentSearchResponse, error) {
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

	// Check if client is initialized
	if s.client == nil {
		log.Error().Msg("gRPC client not initialized")
		return entities.ComponentSearchResponse{}, fmt.Errorf("gRPC client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	grpcRequest := s.mapper.MapToGrpcSearchRequest(request)

	grpcResponse, err := s.client.SearchComponents(s.createAuthContext(ctx), &grpcRequest)
	if err != nil {
		log.Error().Err(err).Msg("Error calling SCANOSS component search API")
		return entities.ComponentSearchResponse{}, fmt.Errorf("API call failed: %w", err)
	}

	response := s.mapper.MapFromGrpcSearchResponse(*grpcResponse)

	log.Debug().
		Int("componentCount", len(response.Components)).
		Msg("Successfully retrieved component search results")

	return response, nil
}

func (s *ScanossApiServiceGrpcImpl) initializeGrpcClient() error {
	if s.endpoint == "" {
		return fmt.Errorf("SCANOSS API endpoint not configured")
	}

	if s.apiKey == "" {
		return fmt.Errorf("SCANOSS API key not configured")
	}

	if s.apiKey != "" && s.endpoint == "" {
		s.endpoint = config.SCANOSS_PREMIUM_API_URL
	}

	config := &tls.Config{
		ServerName: s.endpoint,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := grpc.NewClient(s.endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return fmt.Errorf("failed to connect to SCANOSS API at %s: %w", s.endpoint, err)
	}

	s.conn = conn
	s.client = componentsv2.NewComponentsClient(conn)

	log.Debug().Str("endpoint", s.endpoint).Msg("Successfully initialized SCANOSS gRPC client")
	return nil
}

func (s *ScanossApiServiceGrpcImpl) createAuthContext(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"x-api-key": s.apiKey,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// Close closes the gRPC connection
func (s *ScanossApiServiceGrpcImpl) Close() error {
	if s.conn != nil {
		log.Debug().Msg("Closing SCANOSS gRPC client connection")
		return s.conn.Close()
	}
	return nil
}
