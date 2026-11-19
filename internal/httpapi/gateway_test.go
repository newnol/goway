package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const chatPath = "/v1/chat/completions"

func TestRoutes(t *testing.T) {
	handler := New("test-key")
	tests := []struct {
		name, method, path, token, body, want string
		status                                int
	}{
		{"health", "GET", "/healthz", "", "", `"status":"ok"`, 200},
		{"models", "GET", "/v1/models", "test-key", "", `"id":"goway-demo"`, 200},
		{"missing key", "GET", "/v1/models", "", "", `"code":"invalid_api_key"`, 401},
		{"wrong key", "POST", chatPath, "wrong", "{", `"code":"invalid_api_key"`, 401},
		{"malformed body", "POST", chatPath, "test-key", "{", `"code":"invalid_request_error"`, 400},
		{"large body", "POST", chatPath, "test-key", strings.Repeat(" ", maxBody+1), "exceeds 1 MiB", 413},
		{"unknown model", "POST", chatPath, "test-key", `{"model":"missing","messages":[{"role":"user"}]}`, `"code":"model_not_found"`, 404},
		{"chat", "POST", chatPath, "test-key", `{"model":"goway-demo","messages":[{"role":"user","content":"hi"}]}`, `"finish_reason":"stop"`, 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.token != "" {
				request.Header.Set("Authorization", "Bearer "+tt.token)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != tt.status || !strings.Contains(response.Body.String(), tt.want) {
				t.Fatalf("got status %d and body %q, want status %d and body containing %q", response.Code, response.Body, tt.status, tt.want)
			}
		})
	}
}

func TestChatStream(t *testing.T) {
	server := httptest.NewServer(New("test-key"))
	defer server.Close()
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL+chatPath,
		strings.NewReader(`{"model":"goway-demo","messages":[{"role":"user","content":"hi"}],"stream":true}`))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer test-key")
	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK || response.Header.Get("Content-Type") != "text/event-stream" {
		t.Fatalf("got status %d with content type %q", response.StatusCode, response.Header.Get("Content-Type"))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	frames := strings.Split(string(body), "\n\n")
	if len(frames) != 5 || frames[3] != "data: [DONE]" || frames[4] != "" {
		t.Fatalf("unexpected stream frames: %q", frames)
	}
	var chunks [3]completion
	for i, frame := range frames[:3] {
		data, ok := strings.CutPrefix(frame, "data: ")
		if !ok {
			t.Fatalf("missing SSE data prefix: %q", frame)
		}
		chunk := &chunks[i]
		if err := json.Unmarshal([]byte(data), chunk); err != nil {
			t.Fatal(err)
		}
		if chunk.ID == "" || chunk.ID != chunks[0].ID || chunk.Object != "chat.completion.chunk" || chunk.Model != demoModel || len(chunk.Choices) != 1 || chunk.Choices[0].Delta == nil {
			t.Fatalf("invalid chat chunk: %s", data)
		}
	}
	first, second, final := chunks[0].Choices[0], chunks[1].Choices[0], chunks[2].Choices[0]
	if first.Delta.Role != "assistant" || second.Delta.Content != demoReply || final.FinishReason == nil || *final.FinishReason != "stop" {
		t.Fatalf("unexpected stream body: %s", body)
	}
}
