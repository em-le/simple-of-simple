package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	// Reset users map for test isolation
	users = make(map[string]string)

	// Successful register
	body := bytes.NewBufferString(`{"username":"testuser","password":"testpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Register with existing user
	body = bytes.NewBufferString(`{"username":"testuser","password":"testpass"}`)
	req = httptest.NewRequest(http.MethodPost, "/register", body)
	w = httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}

	// Invalid method
	req = httptest.NewRequest(http.MethodGet, "/register", nil)
	w = httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	users = map[string]string{"testuser": "testpass"}

	// Successful login
	body := bytes.NewBufferString(`{"username":"testuser","password":"testpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	w := httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Wrong password
	body = bytes.NewBufferString(`{"username":"testuser","password":"wrong"}`)
	req = httptest.NewRequest(http.MethodPost, "/login", body)
	w = httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// User not found
	body = bytes.NewBufferString(`{"username":"nouser","password":"testpass"}`)
	req = httptest.NewRequest(http.MethodPost, "/login", body)
	w = httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// Invalid method
	req = httptest.NewRequest(http.MethodGet, "/login", nil)
	w = httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestLogoutHandler(t *testing.T) {
	sessions = make(map[string]string)
	sessions["testuser_session"] = "testuser"

	// Successful logout
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "testuser_session"})
	w := httptest.NewRecorder()
	logoutHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// No session cookie
	req = httptest.NewRequest(http.MethodPost, "/logout", nil)
	w = httptest.NewRecorder()
	logoutHandler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// Invalid method
	req = httptest.NewRequest(http.MethodGet, "/logout", nil)
	w = httptest.NewRecorder()
	logoutHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
