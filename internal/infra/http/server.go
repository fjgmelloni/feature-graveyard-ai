package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"feature-graveyard-ai/internal/application/graveyard"
	"feature-graveyard-ai/internal/domain/feature"
)

type Server struct {
	service graveyard.Service
	mux     *http.ServeMux
}

func NewServer(service graveyard.Service, staticFS http.FileSystem) Server {
	server := Server{
		service: service,
		mux:     http.NewServeMux(),
	}

	server.routes(staticFS)
	return server
}

func (s Server) Handler() http.Handler {
	return s.mux
}

func (s Server) routes(staticFS http.FileSystem) {
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("POST /api/usage-logs", s.ingestUsageLogs)
	s.mux.HandleFunc("GET /api/graveyard/report", s.report)
	s.mux.Handle("/", http.FileServer(staticFS))
}

func (s Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "feature-graveyard-ai"})
}

func (s Server) ingestUsageLogs(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Logs []feature.UsageLogInput `json:"logs"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if len(request.Logs) == 0 {
		writeError(w, http.StatusBadRequest, "logs cannot be empty")
		return
	}

	logs, err := s.service.Ingest(r.Context(), request.Logs)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"ingested": len(logs)})
}

func (s Server) report(w http.ResponseWriter, r *http.Request) {
	windowDays := 180
	if rawWindow := r.URL.Query().Get("windowDays"); rawWindow != "" {
		parsed, err := strconv.Atoi(rawWindow)
		if err != nil {
			writeError(w, http.StatusBadRequest, "windowDays must be a number")
			return
		}
		windowDays = parsed
	}

	report, err := s.service.Report(r.Context(), windowDays)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, report)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}
