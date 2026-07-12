package prometheus

import (
	"context"
	"net/http"
	"strconv"
	"time"

	client "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Registry struct {
	inner        *client.Registry
	httpRequests *client.CounterVec
	httpDuration *client.HistogramVec
	grpcRequests *client.CounterVec
}

func NewRegistry() *Registry {
	registry := &Registry{
		inner: client.NewRegistry(),
		httpRequests: client.NewCounterVec(client.CounterOpts{
			Name: "http_requests_total", Help: "Total HTTP requests.",
		}, []string{"method", "path", "status"}),
		httpDuration: client.NewHistogramVec(client.HistogramOpts{
			Name: "http_request_duration_seconds", Help: "HTTP request duration.",
		}, []string{"method", "path"}),
		grpcRequests: client.NewCounterVec(client.CounterOpts{
			Name: "grpc_requests_total", Help: "Total gRPC requests.",
		}, []string{"method", "code"}),
	}
	registry.inner.MustRegister(registry.httpRequests, registry.httpDuration, registry.grpcRequests)
	return registry
}

func (r *Registry) Registerer() client.Registerer { return r.inner }
func (r *Registry) Handler() http.Handler {
	return promhttp.HandlerFor(r.inner, promhttp.HandlerOpts{})
}

func (r *Registry) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		started := time.Now()
		wrapped := &statusWriter{ResponseWriter: writer, status: http.StatusOK}
		next.ServeHTTP(wrapped, request)
		r.httpRequests.WithLabelValues(request.Method, request.URL.Path, strconv.Itoa(wrapped.status)).Inc()
		r.httpDuration.WithLabelValues(request.Method, request.URL.Path).Observe(time.Since(started).Seconds())
	})
}

func (r *Registry) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		response, err := handler(ctx, req)
		r.grpcRequests.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()
		return response, err
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
