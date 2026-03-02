package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"docs-aggregation-service/internal/metrics"
)

type HTTPServer struct {
	aggregateUsecase AggregateUsecase
	getStatusUsecase GetStatusUsecase
	getResultUsecase GetResultUsecase
	metrics          *metrics.Metrics
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(statusCode int) {
	sr.status = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}

func NewHTTPServer(
	aggregateUsecase AggregateUsecase,
	getStatusUsecase GetStatusUsecase,
	getResultUsecase GetResultUsecase,
	metrics *metrics.Metrics,
) *http.Server {
	httpServer := &HTTPServer{
		aggregateUsecase: aggregateUsecase,
		getStatusUsecase: getStatusUsecase,
		getResultUsecase: getResultUsecase,
		metrics:          metrics,
	}

	mux := http.NewServeMux()
	mux.Handle("/aggregation", httpServer.metricsMiddleware(http.HandlerFunc(httpServer.handleAggregation)))
	mux.Handle("/aggregation/status", httpServer.metricsMiddleware(http.HandlerFunc(httpServer.handleStatus)))
	mux.Handle("/aggregation/result", httpServer.metricsMiddleware(http.HandlerFunc(httpServer.handleResult)))
	mux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))

	return &http.Server{Handler: mux}
}

func (s *HTTPServer) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := &statusRecorder{
			ResponseWriter: w,
			status:         200,
		}
		next.ServeHTTP(sr, r)
		s.metrics.HTTPRequestsTotal.WithLabelValues(r.Pattern, fmt.Sprint(sr.status)).Inc()
	})
}

func (s *HTTPServer) handleAggregation(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HTTPServer] /aggregation from %s", r.RemoteAddr)

	if r.Method != http.MethodPost {
		log.Printf("[HTTPServer] /aggregation: (405) invalid method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	startDateString := r.URL.Query().Get("startDate")
	endDateString := r.URL.Query().Get("endDate")
	if startDateString == "" || endDateString == "" {
		log.Printf("[HTTPServer] /aggregation: (400) startDate or endDate is not provided")
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "startDate and endDate are required"})
		return
	}
	startDate, err := time.Parse(time.DateTime, startDateString)
	if err != nil {
		log.Printf("[HTTPServer] /aggregation: (400) invalid startDate: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid startDate"})
		return
	}
	endDate, err := time.Parse(time.DateTime, endDateString)
	if err != nil {
		log.Printf("[HTTPServer] /aggregation: (400) invalid endDate: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid endDate"})
		return
	}
	if startDate.Compare(endDate) == 1 {
		log.Printf("[HTTPServer] /aggregation: (400) startDate greater than endDate")
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "startDate must be equal or earlier than endDate"})
		return
	}

	taskID, err := s.aggregateUsecase.Run(startDate, endDate)
	if err != nil {
		log.Printf("[HTTPServer] /aggregation: (500) internal AggregationUsecase error")
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "starting aggregation failed"})
		return
	}
	log.Printf("[HTTPServer] /aggregation: (202) task started (ID %s...)", taskID[:6])
	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{"taskID": taskID})
}

func (s *HTTPServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[HTTPServer] /aggregation/status: (405) invalid method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("taskID")
	tasks, err := s.getStatusUsecase.Run(taskID)
	if err != nil {
		log.Printf("[HTTPServer] /aggregation/status: (500) internal GetStatusUsecase error")
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "getting statuses failed"})
		return
	}
	log.Printf("[HTTPServer] /aggregation/status: (200) getting statuses completed")
	w.WriteHeader(http.StatusOK)
	writeJSON(w, tasks)
}

func (s *HTTPServer) handleResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[HTTPServer] /aggregation/result: (405) invalid method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("taskID")
	if taskID == "" {
		log.Printf("[HTTPServer] /aggregation/result: (400) taskID is not provided")
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "taskID is required"})
		return
	}
	filepath, err := s.getResultUsecase.Run(taskID)
	if err != nil {
		if err.Error() == "not completed" {
			log.Printf("[HTTPServer] /aggregation/result: (400) task is not completed (ID %s...)", taskID[:6])
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, map[string]string{"error": "task is not completed"})
			return
		}
		log.Printf("[HTTPServer] /aggregation/result: (500) internal GetResultUsecase error (ID %s...)", taskID[:6])
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "getting result failed"})
		return
	}
	http.ServeFile(w, r, filepath)
	log.Printf("[HTTPServer] /aggregation/result: (200) getting result completed (ID %s...)", taskID[:6])
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
