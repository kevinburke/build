// +build appenginevm

package devapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetTokenMethodNotAllowed(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/setToken", nil)
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Errorf("GET /setToken: got %d, want 405", w.Code)
	}
}
