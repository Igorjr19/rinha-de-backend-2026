package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"

	"Igorjr19/rinha-de-backend-2026/internal/fraud"
)

const (
	maxBodyBytes    = 8 << 10
	maxPooledBuffer = 64 << 10
)

var (
	requestPool = sync.Pool{New: func() any { return new(fraud.Request) }}
	bodyPool    = sync.Pool{New: func() any { b := make([]byte, 0, 2048); return &b }}
	respPool    = sync.Pool{New: func() any { b := make([]byte, 0, 80); return &b }}
)

func handleFraudScore(scorer *fraud.Scorer, w http.ResponseWriter, r *http.Request) {
	bodyPtr := bodyPool.Get().(*[]byte)
	body := (*bodyPtr)[:0]
	defer func() {
		if cap(body) <= maxPooledBuffer {
			*bodyPtr = body[:0]
			bodyPool.Put(bodyPtr)
		}
	}()

	if cl := r.ContentLength; cl > 0 {
		if cl > maxBodyBytes {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		if int64(cap(body)) < cl {
			body = make([]byte, cl)
		} else {
			body = body[:cl]
		}
		if _, err := io.ReadFull(r.Body, body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		var err error
		body, err = io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	req := requestPool.Get().(*fraud.Request)
	*req = fraud.Request{}
	defer requestPool.Put(req)

	if err := json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := scorer.Score(req)
	writeFraudResponse(w, resp)
}

func writeFraudResponse(w http.ResponseWriter, resp fraud.Response) {
	outPtr := respPool.Get().(*[]byte)
	out := (*outPtr)[:0]

	out = append(out, `{"approved":`...)
	out = strconv.AppendBool(out, resp.Approved)
	out = append(out, `,"fraud_score":`...)
	out = strconv.AppendFloat(out, resp.FraudScore, 'f', -1, 64)
	out = append(out, '}')

	h := w.Header()
	h["Content-Type"] = []string{"application/json"}
	h["Content-Length"] = []string{strconv.Itoa(len(out))}
	w.Write(out)

	*outPtr = out[:0]
	respPool.Put(outPtr)
}
