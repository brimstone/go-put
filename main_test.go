package main

import (
	"net/http/httptest"
	"testing"

	"github.com/brimstone/go-saverequest"
)

func TestHandleData(t *testing.T) {
	saverequest.TestRequestFiles(t, ".", handleData)
}

func TestAutoIndex(t *testing.T) {
	// need to reset root
	t.Log("Trying empty index.")
	req, _ := saverequest.FakeRequest("GET", "/data/", map[string]string{}, "")
	w := httptest.NewRecorder()
	handleData(w, req)
	if w.Code != 404 && w.Body.String() != "" {
		t.Errorf("Unable to get empty request")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	} else {
		t.Log("Got generated index")
	}

	t.Log("Trying to put an index file")
	req, _ = saverequest.FakeRequest("POST", "/data/.index", map[string]string{"Content-Type": "text/html"}, "This is the index")
	w = httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "" {
		t.Errorf("Unable to save new index request")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	}

	req, _ = saverequest.FakeRequest("GET", "/data/", map[string]string{"Content-Type": "text/html"}, "")
	w = httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "This is the index" {
		t.Errorf("Unable to get new index request")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	}

	req, _ = saverequest.FakeRequest("GET", "/data/.index", map[string]string{"Content-Type": "text/html"}, "")
	w = httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "This is the index" {
		t.Errorf("Unable to generated index")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	}

}

func TestIndex(t *testing.T) {
	// need to reset root
	t.Log("Trying empty index.")
	req, _ := saverequest.FakeRequest("GET", "/data/", map[string]string{}, "")
	w := httptest.NewRecorder()
	handleData(w, req)
	if w.Code != 404 && w.Body.String() != "" {
		t.Errorf("Unable to get empty request")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	} else {
		t.Log("Got generated index")
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

func TestTree(t *testing.T) {
	// need to reset root
	root = dataItem{}

	t.Log("Trying to put an index file")
	req, _ := saverequest.FakeRequest("POST", "/data/.index", map[string]string{"Content-Type": "text/html"}, "This is the index")
	w := httptest.NewRecorder()
	handleData(w, req)
	if w.Body.String() != "" {
		t.Errorf("Unable to save new index request")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	}

	req, _ = saverequest.FakeRequest("GET", "/tree/", map[string]string{}, "")
	w = httptest.NewRecorder()
	handleTree(w, req)
	t.Log("tree:\n", w.Body.String())

}
