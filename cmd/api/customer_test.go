// customer tests
package main

import (
	"net/http"
	"strings"
	"testing"
	"io"
	"net/http/httptest"
	"log"
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func NewApplication() (*Application, error) {
	var cfg Config

	cfg.port = 4000
	cfg.env = "development"
	cfg.db.dsn = "postgres://postgres:password@localhost/savanah?sslmode=disable"

	logger := log.New(os.Stdout, "", log.LstdFlags)

	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	logger.Printf("database connection pool established")

	app := &Application{
		config: cfg,
		logger: logger,
		DB:     db,
	}

	return app, nil
}

func newTestApplication(t *testing.T) *Application {
	t.Helper()
	app, err := NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	return app
}

func TestCreateCustomer(t *testing.T) {

	app := newTestApplication(t)

	body := `{"name": "Test User", "email": "test@mail.com", "password": "password123"}`
	r, err := http.NewRequest("POST", "/v1/customers", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	app.SignUpCustomer(w, r)

	rs := w.Result()

	if rs.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201; got %d", rs.StatusCode)
	}

	// Check the response body is as expected.
	expected := `{"id":1,"name":"Test User","email":"test@mail.com"}`

	// Read the response body from the http.Response.
	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response body is what we expect.
	if string(b) != expected {
		t.Errorf("expected body %q; got %q", expected, string(b))
	}
}

func TestGetCustomer(t *testing.T) {
	app := newTestApplication(t)

	r, err := http.NewRequest("GET", "/v1/customers/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	app.handle_getSingleCustomer(w, r)

	rs := w.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status 200; got %d", rs.StatusCode)
	}

	// Check the response body is as expected.
	expected := `{"id":1,"name":"Test User","email":"test@mail.com"}`
	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Errorf("expected body %q; got %q", expected, string(b))
	}
}