package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	logger           *log.Logger
	callbackRequest  func(string)
	callbackResponse func(string, int, time.Duration)
)

// Request - структура для работы с запросом
type Request struct {
	w         http.ResponseWriter
	r         *http.Request
	beginTime time.Time
	body      []byte
	route     string
	requestID string
}

// New - функция создает новый запрос
func New(w http.ResponseWriter, r *http.Request) (request *Request) {
	if callbackRequest == nil {
		callbackRequest = dummyCallbackRequest
	}
	if callbackResponse == nil {
		callbackResponse = dummyCallbackResponse
	}
	request = &Request{
		w:         w,
		r:         r,
		beginTime: time.Now(),
	}
	var err error
	if r.Body != nil {
		defer r.Body.Close()
		if request.body, err = ioutil.ReadAll(r.Body); err != nil {
			request.Log().Errorf("Unable to read request body: %s", err)
		} else {
			r.Body = ioutil.NopCloser(bytes.NewReader(request.body))
		}
	}
	l := request.Log()
	if request.body != nil && len(request.body) > 0 && len(request.body) < (1<<20) {
		l = l.WithField("body", string(request.body))
	}
	muxRoute := mux.CurrentRoute(r)
	if muxRoute != nil {
		request.route, _ = muxRoute.GetPathTemplate()
	}

	// request ID
	if value := r.Context().Value(ContextKeyRequestID); value != nil {
		request.requestID = value.(string)
	}
	request.Log().Debugf("Used request ID: %s", request.requestID)

	go callbackRequest(request.route)

	l.Debug("Request")
	return
}

// Setup - функция устанавливает логгер и коллбэки
func Setup(
	l *log.Logger,
	req func(string),
	resp func(string, int, time.Duration),
) {
	logger = l
	callbackRequest = req
	callbackResponse = resp
}

// Log - функция возвращает обогащенный logger для запроса
func (r *Request) Log() *log.Entry {
	if logger == nil {
		logger = log.New()
	}
	entry := logger.
		WithField("method", r.r.Method).
		WithField("host", r.r.Host).
		WithField("proto", r.r.Proto).
		WithField("remote_addr", r.r.RemoteAddr).
		WithField("request_uri", r.r.RequestURI).
		WithField("route", r.route).
		WithField("duration", time.Now().Sub(r.beginTime))

	if r.requestID != "" {
		entry.WithField("request_id", r.requestID)
	}

	return entry
}

// FinishOK функция завершает запрос удачно с кодом 200
func (r *Request) FinishOK(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusOK).
		Infof("Response: %s", fmt.Sprintf(msg, args...))
	r.finish(http.StatusOK, msg, args...)
}

// FinishBadRequest функция завершает запрос неудачно с кодом 400
func (r *Request) FinishBadRequest(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusBadRequest).
		Warnf("Response: %s", fmt.Sprintf(msg, args...))
	r.finish(http.StatusBadRequest, msg, args...)
}

// FinishError функция завершает запрос неудачно с кодом 500
func (r *Request) FinishError(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusOK).
		Errorf("Response: %s", fmt.Sprintf(msg, args...))
	r.finish(http.StatusInternalServerError, msg, args...)
}

// FinishOKJSON функция завершает запрос с кодом 200 и объектом для JSON
func (r *Request) FinishOKJSON(i interface{}) {
	r.FinishJSON(http.StatusOK, i)
}

// FinishJSON функция завершает запрос с произвольным кодом и объектом для JSON
func (r *Request) FinishJSON(code int, i interface{}) {
	data, err := json.Marshal(i)
	if err != nil {
		r.Log().Errorf("Unable to marshal response data: %s", err)
		r.FinishError("Unable to marshal response data: %s", err)
		return
	}

	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(code)
	if _, err := r.w.Write(data); err != nil {
		r.Log().Warnf("Unable to write data: %s", err)
		return
	}
	ll := r.Log().
		WithField("status", code)
	if len(data) < 1048576 {
		ll = ll.WithField("body", string(data))
	}
	ll.Infof("Response")
	go callbackResponse(r.route, http.StatusOK, time.Since(r.beginTime))
}

// Finish функция завершает запрос с введенным кодом
func (r *Request) Finish(code int, msg string, args ...interface{}) {
	r.Log().
		WithField("status", code).
		Infof("Response: %s", fmt.Sprintf(msg, args...))
	r.finish(code, msg, args...)
}

// FinishNoContent функция завершает запрос с кодом 204
func (r *Request) FinishNoContent() {
	r.Log().
		WithField("status", http.StatusNoContent).
		Infof("Response no content")
	r.w.WriteHeader(http.StatusNoContent)
	go callbackResponse(r.route, http.StatusNoContent, time.Since(r.beginTime))
}

// FinishFile - функция завершает запрос с указанным кодом,
// передавая данные байты, как файл filename с указанным contentType
func (r *Request) FinishFile(code int, filename, contentType string, data []byte) {
	r.w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	r.w.Header().Add("Content-Type", contentType)
	r.w.WriteHeader(code)
	if _, err := r.w.Write(data); err != nil {
		r.Log().Warnf("Unable to write data as file: %s", err)
		return
	}
	ll := r.Log().
		WithField("status", code)
	ll.Infof("Response")
	go callbackResponse(r.route, code, time.Since(r.beginTime))
}

// GetVar функция возвращает переменную пути по имени
func (r *Request) GetVar(name string) string {
	return mux.Vars(r.r)[name]
}

// GetBody - функция извлекает
func (r *Request) GetBody() []byte {
	return r.body
}

// Query - функция возвращает query-параметры
func (r *Request) Query() url.Values {
	return r.r.URL.Query()
}

// QueryValue - функция возвращает по имени аргумент запроса
func (r *Request) QueryValue(name string) Value {
	return Value(r.Query().Get(name))
}

// VarsValue - функция возвращает по имени переменную пути
func (r *Request) VarsValue(name string) Value {
	return Value(r.GetVar(name))
}

func (r *Request) finish(code int, msg string, args ...interface{}) {
	r.w.WriteHeader(code)
	buf := bytes.NewBufferString(fmt.Sprintf(msg, args...))
	r.w.Write(buf.Bytes())
	go callbackResponse(r.route, code, time.Since(r.beginTime))
}

func dummyCallbackRequest(s string) {
}

func dummyCallbackResponse(s string, code int, duration time.Duration) {
}
