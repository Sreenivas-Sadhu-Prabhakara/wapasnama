package backend

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// memStore is an in-memory Store for handler tests (no DB).
type memStore struct{ items []Record }

func (m *memStore) Save(r Record) (Record, error) {
	r.ID = int64(len(m.items) + 1)
	m.items = append([]Record{r}, m.items...)
	return r, nil
}
func (m *memStore) List(limit int) ([]Record, error) { return m.items, nil }

func TestLogAndSummary(t *testing.T) {
	srv := NewServer(&memStore{})
	post := func(body string) int {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/log", strings.NewReader(body)))
		return rec.Code
	}
	if c := post(`{"supplier":"ACME","item":"soap","kind":"short","amount":300,"status":"open"}`); c != http.StatusCreated {
		t.Fatalf("log status %d", c)
	}
	if c := post(`{"supplier":"ACME","item":"oil","kind":"damaged","amount":120,"status":"received"}`); c != http.StatusCreated {
		t.Fatalf("log status %d", c)
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/summary", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"openAmount":300`) {
		t.Fatalf("summary=%s", rec.Body.String())
	}
}

func TestLog_ValidationError(t *testing.T) {
	srv := NewServer(&memStore{})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/log", strings.NewReader(`{"kind":"short","status":"open"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d want 400", rec.Code)
	}
}
