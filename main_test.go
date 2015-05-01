package main

import (
	"net/http/httptest"
	"testing"

	"github.com/brimstone/go-saverequest"
)

func TestHandleData(t *testing.T) {
	saverequest.TestRequestFiles(t, ".", handleData)
}

func TestIndex(t *testing.T) {
	req, _ := saverequest.FakeRequest("GET", "/data/", map[string]string{}, "")
	w := httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "[]" {
		t.Errorf("Unable to get the index request")
		return
	}

	req, _ = saverequest.FakeRequest("POST", "/data/.index", map[string]string{"Content-Type": "text/html"}, "This is the index")
	w = httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "" {
		t.Errorf("Unable to save new index request")
		return
	}

	req, _ = saverequest.FakeRequest("GET", "/data/", map[string]string{"Content-Type": "text/html"}, "")
	w = httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "This is the index" {
		t.Errorf("Unable to save new index request")
		return
	}

	req, _ = saverequest.FakeRequest("GET", "/data/.index", map[string]string{"Content-Type": "text/html"}, "")
	w = httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "This is the index" {
		t.Errorf("Unable to save new index request")
		return
	}

}

func Test404(t *testing.T) {
	req, _ := saverequest.FakeRequest("GET", "/data/notfound", map[string]string{}, "")
	w := httptest.NewRecorder()
	handleData(w, req)
	if w.Code != 404 {
		t.Errorf("Expected code 404, got", w.Code)
		return
	}

}
