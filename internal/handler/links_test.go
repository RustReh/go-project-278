package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RustReh/go-project-278/internal/schemas"
)

func TestLinks_Create_WithShortName_201(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := `{"original_url":"https://example.com/long","short_name":"exmpl"}`
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp schemas.LinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ShortName != "exmpl" || resp.ShortURL != "https://short.io/r/exmpl" {
		t.Fatalf("got %+v", resp)
	}
}

func TestLinks_Create_WithoutShortName_201(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := `{"original_url":"https://example.com/auto"}`
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp schemas.LinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ShortName == "" {
		t.Fatal("expected generated short_name")
	}
	if resp.ShortURL != handlerBaseURL+"r/"+resp.ShortName {
		t.Fatalf("short_url: got %q", resp.ShortURL)
	}
}

func TestLinks_Create_Conflict_422(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := `{"original_url":"https://example.com/1","short_name":"dup"}`
	for range 2 {
		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code == http.StatusCreated {
			continue
		}
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status: got %d, body: %s", rec.Code, rec.Body.String())
		}
		var resp schemas.ValidationErrorsResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		if resp.Errors["short_name"] != "short name already in use" {
			t.Fatalf("errors: %#v", resp.Errors)
		}
		return
	}
	t.Fatal("expected conflict on second create")
}

func TestLinks_GetAll_200(t *testing.T) {
	r, _ := setupTestRouter(t)

	create := func(payload string) {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create status %d: %s", rec.Code, rec.Body.String())
		}
	}
	create(`{"original_url":"https://example.com/1","short_name":"aaa"}`)
	create(`{"original_url":"https://example.com/2","short_name":"bbb"}`)

	req := httptest.NewRequest(http.MethodGet, "/api/links?range=[0,10]", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, body: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Range"); got != "links 0-10/2" {
		t.Fatalf("Content-Range: got %q, want links 0-10/2", got)
	}
	var resp []schemas.LinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Fatalf("len: got %d, want 2", len(resp))
	}
}

func TestLinks_GetAll_Pagination(t *testing.T) {
	r, _ := setupTestRouter(t)

	for i := 1; i <= 11; i++ {
		payload := `{"original_url":"https://example.com/` + string(rune('0'+i%10)) + `","short_name":"ln` + string(rune('0'+i)) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create %d: status %d %s", i, rec.Code, rec.Body.String())
		}
	}

	t.Run("first page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/links?range=[0,10]", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: %d", rec.Code)
		}
		if got := rec.Header().Get("Content-Range"); got != "links 0-10/11" {
			t.Fatalf("Content-Range: got %q", got)
		}
		var resp []schemas.LinkResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		if len(resp) != 10 {
			t.Fatalf("len: got %d, want 10", len(resp))
		}
		if resp[0].ID != 1 || resp[9].ID != 10 {
			t.Fatalf("ids: %d-%d", resp[0].ID, resp[9].ID)
		}
	})

	t.Run("second slice", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/links?range=[5,10]", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if got := rec.Header().Get("Content-Range"); got != "links 5-10/11" {
			t.Fatalf("Content-Range: got %q", got)
		}
		var resp []schemas.LinkResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		if len(resp) != 5 {
			t.Fatalf("len: got %d, want 5", len(resp))
		}
		if resp[0].ID != 6 || resp[4].ID != 10 {
			t.Fatalf("ids: %d-%d", resp[0].ID, resp[4].ID)
		}
	})
}

func TestLinks_GetAll_MissingRange_422(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/links", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status: got %d, want 422", rec.Code)
	}
}

func TestLinks_GetByID_200_and_404(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/links",
		bytes.NewBufferString(`{"original_url":"https://example.com/x","short_name":"xxx"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	var created schemas.LinkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/links/1", nil)
	getRec := httptest.NewRecorder()
	r.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status: %d, body: %s", getRec.Code, getRec.Body.String())
	}

	missReq := httptest.NewRequest(http.MethodGet, "/api/links/999", nil)
	missRec := httptest.NewRecorder()
	r.ServeHTTP(missRec, missReq)
	if missRec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", missRec.Code)
	}
}

func TestLinks_Update_200_and_404(t *testing.T) {
	r, _ := setupTestRouter(t)

	postReq := httptest.NewRequest(http.MethodPost, "/api/links",
		bytes.NewBufferString(`{"original_url":"https://example.com/old","short_name":"old"}`))
	postReq.Header.Set("Content-Type", "application/json")
	postRec := httptest.NewRecorder()
	r.ServeHTTP(postRec, postReq)

	putReq := httptest.NewRequest(http.MethodPut, "/api/links/1",
		bytes.NewBufferString(`{"original_url":"https://example.com/new","short_name":"new"}`))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	r.ServeHTTP(putRec, putReq)
	if putRec.Code != http.StatusOK {
		t.Fatalf("update status: %d, body: %s", putRec.Code, putRec.Body.String())
	}

	var updated schemas.LinkResponse
	if err := json.Unmarshal(putRec.Body.Bytes(), &updated); err != nil {
		t.Fatal(err)
	}
	if updated.ShortName != "new" {
		t.Fatalf("got %+v", updated)
	}

	missPut := httptest.NewRequest(http.MethodPut, "/api/links/999",
		bytes.NewBufferString(`{"original_url":"https://example.com/x","short_name":"xxx"}`))
	missPut.Header.Set("Content-Type", "application/json")
	missRec := httptest.NewRecorder()
	r.ServeHTTP(missRec, missPut)
	if missRec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", missRec.Code)
	}
}

func TestLinks_Delete_204_and_404(t *testing.T) {
	r, _ := setupTestRouter(t)

	postReq := httptest.NewRequest(http.MethodPost, "/api/links",
		bytes.NewBufferString(`{"original_url":"https://example.com/d","short_name":"del"}`))
	postReq.Header.Set("Content-Type", "application/json")
	postRec := httptest.NewRecorder()
	r.ServeHTTP(postRec, postReq)

	delReq := httptest.NewRequest(http.MethodDelete, "/api/links/1", nil)
	delRec := httptest.NewRecorder()
	r.ServeHTTP(delRec, delReq)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status: got %d", delRec.Code)
	}
	if delRec.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", delRec.Body.String())
	}

	missDel := httptest.NewRequest(http.MethodDelete, "/api/links/1", nil)
	missRec := httptest.NewRecorder()
	r.ServeHTTP(missRec, missDel)
	if missRec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", missRec.Code)
	}
}
