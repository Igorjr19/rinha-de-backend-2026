package api

import (
	"net/http"

	"Igorjr19/rinha-de-backend-2026/internal/fraud"
)

type router struct {
	scorer *fraud.Scorer
}

func NewRouter(scorer *fraud.Scorer) http.Handler {
	return &router{scorer: scorer}
}

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/fraud-score":
		if r.Method == http.MethodPost {
			handleFraudScore(rt.scorer, w, r)
			return
		}
	case "/ready":
		if r.Method == http.MethodGet {
			if !rt.scorer.Ready() {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
