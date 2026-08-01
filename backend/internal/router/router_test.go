package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-blog/backend/internal/middleware"
)

func TestHealth(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)

	New(Dependencies{}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data.Status != "ok" || response.Message != "success" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestCORSAllowsConfiguredOriginAndPreflight(t *testing.T) {
	engine := New(Dependencies{
		CORSMiddleware: middleware.CORS([]string{"https://blog.example.com"}),
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/health", nil)
	request.Header.Set("Origin", "https://blog.example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
	if recorder.Header().Get("Access-Control-Allow-Origin") != "https://blog.example.com" {
		t.Fatalf("unexpected allow origin header: %q", recorder.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSRejectsUnknownOriginPreflight(t *testing.T) {
	engine := New(Dependencies{
		CORSMiddleware: middleware.CORS([]string{"https://blog.example.com"}),
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/health", nil)
	request.Header.Set("Origin", "https://attacker.example")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
	if recorder.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("unexpected allow origin header: %q", recorder.Header().Get("Access-Control-Allow-Origin"))
	}
}
