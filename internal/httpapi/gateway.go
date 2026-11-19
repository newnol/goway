package httpapi

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	demoModel = "goway-demo"
	demoReply = "Hello from GoWay's mock provider."
	maxBody   = 1 << 20
)

type message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type choice struct {
	Index        int      `json:"index"`
	Message      *message `json:"message,omitempty"`
	Delta        *message `json:"delta,omitempty"`
	FinishReason *string  `json:"finish_reason"`
}

type completion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
}

func New(apiKey string) http.Handler {
	keyHash := sha256.Sum256([]byte(apiKey))
	authenticate := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			scheme, token, _ := strings.Cut(r.Header.Get("Authorization"), " ")
			tokenHash := sha256.Sum256([]byte(token))
			if !strings.EqualFold(scheme, "Bearer") || token == "" || subtle.ConstantTimeCompare(keyHash[:], tokenHash[:]) != 1 {
				w.Header().Set("WWW-Authenticate", "Bearer")
				writeError(w, http.StatusUnauthorized, "invalid_api_key", "Invalid or missing bearer token")
				return
			}
			next(w, r)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/models", authenticate(listModels))
	mux.HandleFunc("POST /v1/chat/completions", authenticate(chat))
	return mux
}

func listModels(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"object": "list",
		"data": []map[string]any{{
			"id": demoModel, "object": "model", "created": 0, "owned_by": "goway",
		}},
	})
}

func chat(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeError(w, http.StatusRequestEntityTooLarge, "invalid_request_error", "Request body exceeds 1 MiB")
		return
	}
	var request struct {
		Model    string `json:"model"`
		Messages []struct {
			Role string `json:"role"`
		} `json:"messages"`
		Stream bool `json:"stream"`
	}
	if err != nil || json.Unmarshal(body, &request) != nil || request.Model == "" || len(request.Messages) == 0 {
		writeError(w, http.StatusBadRequest, "invalid_request_error", "Provide a model and a non-empty messages array")
		return
	}
	for _, msg := range request.Messages {
		if msg.Role == "" {
			writeError(w, http.StatusBadRequest, "invalid_request_error", "Each message requires a role")
			return
		}
	}
	if request.Model != demoModel {
		writeError(w, http.StatusNotFound, "model_not_found", "Unknown model")
		return
	}

	reply := completion{
		ID:      "chatcmpl-" + rand.Text(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   demoModel,
	}
	stop := "stop"
	if request.Stream {
		streamChat(w, reply, []choice{
			{Delta: &message{Role: "assistant"}},
			{Delta: &message{Content: demoReply}},
			{Delta: &message{}, FinishReason: &stop},
		})
		return
	}
	reply.Choices = []choice{{Message: &message{Role: "assistant", Content: demoReply}, FinishReason: &stop}}
	writeJSON(w, http.StatusOK, reply)
}

func streamChat(w http.ResponseWriter, reply completion, chunks []choice) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	controller := http.NewResponseController(w)
	reply.Object = "chat.completion.chunk"
	for _, chunk := range chunks {
		reply.Choices = []choice{chunk}
		data, err := json.Marshal(reply)
		if err != nil {
			return
		}
		if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
			return
		}
		if err := controller.Flush(); err != nil {
			return
		}
	}
	if _, err := io.WriteString(w, "data: [DONE]\n\n"); err == nil {
		_ = controller.Flush()
	}
}

func writeError(w http.ResponseWriter, status int, code, detail string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{
		"message": detail,
		"type":    "invalid_request_error",
		"code":    code,
	}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
