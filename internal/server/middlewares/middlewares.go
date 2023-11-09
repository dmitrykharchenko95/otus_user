package middlewares

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dmitrykharchenko95/otus_user/internal/server/response"
)

const tokenLength = 16

func init() {
	prometheus.MustRegister(requestCounter, responseStatus, requestDuration)

}

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		}, []string{"http_method", "http_endpoint"})

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status",
			Help: "Status of HTTP response.",
		}, []string{"http_method", "http_endpoint", "http_response_status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		}, []string{"http_method", "http_endpoint"})

	Default = []mux.MiddlewareFunc{RequestID, Logging, Metrics}
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if exclude(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		var token = generateToken()
		r.Header.Add("X-Request-ID", token)

		next.ServeHTTP(response.New(w), r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if exclude(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		var reqBody, _ = io.ReadAll(r.Body)
		log.Printf("%s [IN] %s %s %s req: %s\n",
			r.Header.Get("X-Request-ID"),
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			string(reqBody),
		)

		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		var startTime = time.Now()

		next.ServeHTTP(w, r)

		var (
			rw, ok = w.(*response.ResponseWriter)
			resp   string
		)
		if ok {
			resp = rw.GetBody()
		}
		log.Printf("%s [OUT] - %dms %s %s %s resp: %s\n",
			r.Header.Get("X-Request-ID"),
			time.Since(startTime).Milliseconds(),
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			resp,
		)
	})
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if exclude(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		var startTime = time.Now()

		requestCounter.With(prometheus.Labels{
			"http_method":   r.Method,
			"http_endpoint": replaceNumbersWithAsterisk(r.RequestURI),
		}).Inc()

		var status int
		switch rw := w.(type) {
		case *response.ResponseWriter:
			next.ServeHTTP(rw, r)
			status = rw.GetStatus()
		default:
			next.ServeHTTP(w, r)
		}

		responseStatus.With(
			prometheus.Labels{
				"http_method":          r.Method,
				"http_endpoint":        replaceNumbersWithAsterisk(r.RequestURI),
				"http_response_status": fmt.Sprintf("%d", status)},
		).Inc()

		requestDuration.With(
			prometheus.Labels{
				"http_method":   r.Method,
				"http_endpoint": replaceNumbersWithAsterisk(r.RequestURI)},
		).Observe(time.Since(startTime).Seconds())
	})
}

func generateToken() string {
	var token = make([]byte, tokenLength)
	_, _ = rand.Read(token)

	return hex.EncodeToString(token)
}

func exclude(path string) bool {
	return path == "/health" || path == "/metrics"
}

func replaceNumbersWithAsterisk(input string) string {
	return regexp.MustCompile(`/\d+`).ReplaceAllString(input, "/*")
}
