// Copyright 2026 Mano
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"net/http"
)

// SetupRoutes sets up the HTTP routes
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	// API routes
	mux.HandleFunc("/api/health", h.Health)
	mux.HandleFunc("/api/config", h.GetConfig)
	mux.HandleFunc("/api/events", h.GetEvents)
	mux.HandleFunc("/api/events/month", h.GetEventsMonth)
}
