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

package mappers

import (
	"github.com/scanoss/papi/api/componentsv2"
	"github.com/scanoss/scanoss.cc/backend/entities"
)

type ScanossApiMapper interface {
	MapToGrpcSearchRequest(request entities.ComponentSearchRequest) componentsv2.CompSearchRequest
	MapFromGrpcSearchResponse(response componentsv2.CompSearchResponse) entities.ComponentSearchResponse
}

// ScanossApiMapperImpl implements the ScanossApiMapper interface
type ScanossApiMapperImpl struct{}

// NewScanossApiMapper creates a new instance of the SCANOSS API mapper
func NewScanossApiMapper() ScanossApiMapper {
	return &ScanossApiMapperImpl{}
}

// MapToGrpcSearchRequest converts internal ComponentSearchRequest to gRPC request
func (m *ScanossApiMapperImpl) MapToGrpcSearchRequest(request entities.ComponentSearchRequest) componentsv2.CompSearchRequest {
	return componentsv2.CompSearchRequest{
		Search:    request.Search,
		Vendor:    request.Vendor,
		Component: request.Component,
		Package:   request.Package,
		Limit:     request.Limit,
		Offset:    request.Offset,
	}
}

// MapFromGrpcSearchResponse converts gRPC response to internal ComponentSearchResponse
func (m *ScanossApiMapperImpl) MapFromGrpcSearchResponse(response componentsv2.CompSearchResponse) entities.ComponentSearchResponse {
	components := make([]entities.SearchedComponent, len(response.Components))

	for i, comp := range response.Components {
		components[i] = entities.SearchedComponent{
			Component: comp.Component,
			Purl:      comp.Purl,
			URL:       comp.Url,
		}
	}

	result := entities.ComponentSearchResponse{
		Components: components,
		Status: entities.StatusResponse{
			Code:    200,
			Message: "Success",
		},
	}

	if response.Status != nil {
		result.Status = entities.StatusResponse{
			Code:    int32(response.Status.Status),
			Message: response.Status.Message,
		}
	}

	return result
}
