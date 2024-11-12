package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/soerenschneider/hermes/internal/domain"
)

type DummyDispatcher struct {
	not *domain.NotificationRequest
	src *string
}

func (d *DummyDispatcher) Accept(notification domain.NotificationRequest, eventSource string) error {
	d.not = &notification
	d.src = &eventSource
	return nil
}

func TestHttpServer_notifyHandler_wrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/notify", getValidReqest(t))
	w := httptest.NewRecorder()

	server := HttpServer{
		dispatcher: &DummyDispatcher{},
	}

	server.notifyHandler(w, req)
	res := w.Result()
	expectedStatus := http.StatusMethodNotAllowed
	defer res.Body.Close()
	if expectedStatus != res.StatusCode {
		t.Errorf("expected %d got %d", expectedStatus, res.StatusCode)
	}
}

func TestHttpServer_notifyHandler_emptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/notify", nil)
	w := httptest.NewRecorder()

	server := HttpServer{
		dispatcher: &DummyDispatcher{},
	}

	server.notifyHandler(w, req)
	res := w.Result()
	expectedStatus := http.StatusBadRequest
	defer res.Body.Close()
	if expectedStatus != res.StatusCode {
		t.Errorf("expected %d got %d", expectedStatus, res.StatusCode)
	}
}

func TestHttpServer_notifyHandler_garbageData(t *testing.T) {
	body := strings.NewReader("garbage")
	req := httptest.NewRequest(http.MethodPost, "/notify", body)
	w := httptest.NewRecorder()

	server := HttpServer{
		dispatcher: &DummyDispatcher{},
	}

	server.notifyHandler(w, req)
	res := w.Result()
	expectedStatus := http.StatusBadRequest
	defer res.Body.Close()
	if expectedStatus != res.StatusCode {
		t.Errorf("expected %d got %d", expectedStatus, res.StatusCode)
	}
}

func TestHttpServer_notifyHandler_ok_statuscode(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/notify", getValidReqest(t))
	w := httptest.NewRecorder()

	server := HttpServer{
		dispatcher: &DummyDispatcher{},
	}

	server.notifyHandler(w, req)
	res := w.Result()
	expectedStatus := http.StatusOK
	defer res.Body.Close()
	if expectedStatus != res.StatusCode {
		t.Errorf("expected %d got %d", expectedStatus, res.StatusCode)
	}
}

func TestHttpServer_notifyHandler_ok(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/notify", getValidReqest(t))
	w := httptest.NewRecorder()

	dispatcher := &DummyDispatcher{}
	server := HttpServer{
		dispatcher: dispatcher,
	}

	server.notifyHandler(w, req)
	res := w.Result()
	expected := domain.NotificationRequest{
		ServiceId: serviceId,
		Subject:   subject,
		Message:   message,
	}

	defer res.Body.Close()

	if expected != *dispatcher.not {
		t.Errorf("expected %v got %v", expected, *dispatcher.not)
	}
}

const (
	serviceId = "example"
	subject   = "example subject"
	message   = "example message"
)

func getValidReqest(t *testing.T) io.Reader {
	req := domain.NotificationRequest{
		ServiceId: serviceId,
		Subject:   subject,
		Message:   message,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Error("can not marshal json, that's a bug", err)
	}

	return bytes.NewReader(data)
}
