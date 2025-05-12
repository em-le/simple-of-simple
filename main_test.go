package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}
}

func clearUsersTable(t *testing.T) {
	db.Exec("DELETE FROM users")
}

func TestRegisterHandler(t *testing.T) {
	setupTestDB(t)
	clearUsersTable(t)

	body := bytes.NewBufferString(`{"firstName":"Test","lastName":"User","email":"testuser@example.com","password":"testpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body = bytes.NewBufferString(`{"firstName":"Test","lastName":"User","email":"testuser@example.com","password":"testpass"}`)
	req = httptest.NewRequest(http.MethodPost, "/register", body)
	w = httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/register", nil)
	w = httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	setupTestDB(t)
	clearUsersTable(t)
	// Create user in DB
	db.Create(&User{FirstName: "Test", LastName: "User", Email: "testuser@example.com", Password: "testpass"})

	body := bytes.NewBufferString(`{"email":"testuser@example.com","password":"testpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	w := httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body = bytes.NewBufferString(`{"email":"testuser@example.com","password":"wrong"}`)
	req = httptest.NewRequest(http.MethodPost, "/login", body)
	w = httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	body = bytes.NewBufferString(`{"email":"nouser@example.com","password":"testpass"}`)
	req = httptest.NewRequest(http.MethodPost, "/login", body)
	w = httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/login", nil)
	w = httptest.NewRecorder()
	loginHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestLogoutHandler(t *testing.T) {
	sessions = make(map[string]string)
	sessions["testuser@example.com_session"] = "testuser@example.com"

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "testuser@example.com_session"})
	w := httptest.NewRecorder()
	logoutHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/logout", nil)
	w = httptest.NewRecorder()
	logoutHandler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/logout", nil)
	w = httptest.NewRecorder()
	logoutHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
