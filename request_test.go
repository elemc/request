package request_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/elemc/request"
	"github.com/gorilla/mux"
)

type mockWriter struct{}

func (mw *mockWriter) Header() http.Header {
	res := make(http.Header)
	return res
}

func (mw *mockWriter) Write(data []byte) (int, error) {
	//fmt.Printf("Buffer: %s\n", data)
	return len(data), nil
}

func (mw *mockWriter) WriteHeader(statusCode int) {
	//fmt.Printf("Status code: %d\n", statusCode)
}

func testRequest() (err error) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	logrus.SetLevel(logrus.ErrorLevel)
	request.Setup(logrus.StandardLogger(), nil, nil)

	router := mux.NewRouter()
	router.HandleFunc("/test/uri", handler)

	u, err := url.Parse("http://localhost/test/uri")
	if err != nil {
		return
	}
	someBody := "Some body"
	w := &mockWriter{}
	buf := ioutil.NopCloser(bytes.NewBufferString(someBody))
	r := &http.Request{
		Body:          buf,
		ContentLength: int64(len(someBody)),
		Host:          "localhost",
		Method:        "GET",
		Proto:         "HTTP",
		ProtoMajor:    1,
		ProtoMinor:    1,
		RemoteAddr:    "localhost:12345",
		RequestURI:    "/test/uri",
		URL:           u,
	}
	apiReq := request.New(w, r)
	body := apiReq.GetBody()
	if string(body) != someBody {
		err = fmt.Errorf("unexpected body: %s, expected: %s", string(body), someBody)
	}
	apiReq.FinishNoContent()

	return
}

func TestRequest(t *testing.T) {
	if err := testRequest(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := testRequest(); err != nil {
			b.Fatal(err)
		}
	}
}
