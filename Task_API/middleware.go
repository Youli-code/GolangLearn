package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		dur := time.Since(start)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, dur)
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				writeErr(w, http.StatusInternalServerError, "internal error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CORS(origins string) Middleware {
	allowedOrigins := parseOrigins(origins)
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	allowedMethods := "GET,POST,PUT,DELETE,OPTIONS"
	allowedHeaders := "Content-Type,Authorization"
	exposeHeaders := "Content-Type"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			origin := r.Header.Get("Origin")
			if origin != "" && (allowAll || containsOrigin(allowedOrigins, origin)) {
				w.Header().Set("Access-Control-Allow-Origin", allowAllIfWildcard(origin, allowAll))
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
				w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
				w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseOrigins(csv string) []string {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return []string{"*"}
	}
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func containsOrigin(list []string, origin string) bool {
	for _, o := range list {
		if o == origin {
			return true
		}
	}
	return false
}

func allowAllIfWildcard(origin string, allowAll bool) string {
	if allowAll {
		return "*"
	}
	return origin
}

func CORSFromEnv() Middleware {
	return CORS(os.Getenv("TASK_API_CORS_ORIGINS"))
}
