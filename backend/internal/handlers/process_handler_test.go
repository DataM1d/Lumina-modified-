package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleAnalyze_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &ProcessHandler{Svc: nil}
	r := gin.New()
	r.POST("/process", h.HandleAnalyze)

	t.Run("Reject Invalid URL Format", func(t *testing.T) {
		body := map[string]string{"url": "not-a-url"}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Reject Missing Field", func(t *testing.T) {
		body := map[string]string{"invalid_key": "https://google.com"}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Handle Nil Service Gracefully", func(t *testing.T) {
		body := map[string]string{"url": "https://example.com"}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	})
}
