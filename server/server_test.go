package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jonasfj/statsum/payload"

	"gopkg.in/dgrijalva/jwt-go.v2"
)

func nilOrPanic(err error, a ...interface{}) {
	if err != nil {
		panic(fmt.Sprintln(append([]interface{}{err}, a...)...))
	}
}

func assert(condition bool, a ...interface{}) {
	if !condition {
		panic(fmt.Sprintln(a...))
	}
}

func doTestRequest(r *http.Request) *httptest.ResponseRecorder {
	if r == nil {
		panic("Failed to provide request, probably error in making it!")
	}
	statsum, err := New(Config{JwtSecret: []byte("secret")})
	nilOrPanic(err)
	w := httptest.NewRecorder()
	statsum.handler(w, r)
	return w
}

var testBody = payload.Payload{
	CountMetrics: []payload.CountMetric{
		{Key: "test-count", Value: 42},
	},
	ValueMetrics: []payload.ValueMetric{
		{Key: "test-measure", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},
	},
}

func jsonBytes(p payload.Payload) []byte {
	b, err := p.MarshalJSON()
	nilOrPanic(err)
	return b
}

func msgpBytes(p payload.Payload) []byte {
	b, err := p.MarshalMsg([]byte{})
	nilOrPanic(err)
	return b
}

func jsonBody(p payload.Payload) io.Reader {
	b, err := p.MarshalJSON()
	nilOrPanic(err)
	return bytes.NewReader(b)
}

func msgpBody(p payload.Payload) io.Reader {
	b, err := p.MarshalMsg([]byte{})
	nilOrPanic(err)
	return bytes.NewReader(b)
}

func auth(project string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["project"] = project
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("secret"))
	nilOrPanic(err, "Failed to sign token")
	return "Bearer " + t
}

func TestHandleCorrectPrefix(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/prefix", jsonBody(testBody))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", auth("prefix"))
	w := doTestRequest(r)
	assert(w.Code == http.StatusOK, "/v1/project/prefix failed")
}

func TestHandleCorrectPrefixSpecialChars(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/Pre-f_ix9", jsonBody(testBody))
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("Authorization", auth("Pre-f_ix9"))
	w := doTestRequest(r)
	assert(w.Code == http.StatusOK, "/v1/project/prefix failed")
}

func TestHandleMsgPack(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/prefix", msgpBody(testBody))
	r.Header.Set("Content-Type", "application/msgpack")
	r.Header.Set("Authorization", auth("prefix"))
	w := doTestRequest(r)
	assert(w.Code == http.StatusOK, "/v1/project/prefix failed")
}

func TestHandleInvalidJSON(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/prefix", strings.NewReader(`{
    wrong json
  }`))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", auth("prefix"))
	w := doTestRequest(r)
	assert(w.Code == http.StatusBadRequest, "invalid json didn't fail")
}

func TestHandleInvalidContentType(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/prefix", jsonBody(testBody))
	r.Header.Set("Content-Type", "something/wrong")
	r.Header.Set("Authorization", auth("prefix"))
	w := doTestRequest(r)
	assert(w.Code == http.StatusUnsupportedMediaType, "invalid content-type didn't fail")
}

func TestHandleInvalidPath(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/lala/lala", jsonBody(testBody))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", auth("lala/lala"))
	w := doTestRequest(r)
	assert(w.Code == http.StatusNotFound, "/lala/lala didn't fail")
}

func TestHandleMissingPrefix(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/", jsonBody(testBody))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", auth(""))
	w := doTestRequest(r)
	assert(w.Code == http.StatusNotFound, "/v1/project/ didn't fail")
}

func TestHandleIllegalPrefix(t *testing.T) {
	r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/d/d", jsonBody(testBody))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", auth("d/d"))
	w := doTestRequest(r)
	assert(w.Code == http.StatusBadRequest, "/v1/project/d/d didn't fail")
}

func TestParallel(t *testing.T) {
	statsum, err := New(Config{JwtSecret: []byte("secret")})
	nilOrPanic(err, "Failed to create statsum")
	authorization := auth("p")
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/p", jsonBody(payload.Payload{
			CountMetrics: []payload.CountMetric{
				{Key: "test-count-1", Value: 1},
			},
		}))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", authorization)
		w := httptest.NewRecorder()
		statsum.handler(w, r)
		assert(w.Code == http.StatusOK, "Failed to send request 1")
		wg.Done()
	}()
	go func() {
		r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/p", jsonBody(payload.Payload{
			CountMetrics: []payload.CountMetric{
				{Key: "test-count-1", Value: 1},
			},
		}))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", authorization)
		w := httptest.NewRecorder()
		statsum.handler(w, r)
		assert(w.Code == http.StatusOK, "Failed to send request 1")
		wg.Done()
	}()
	wg.Wait()
}

var pc1 = payload.Payload{
	CountMetrics: []payload.CountMetric{
		{Key: "test-count-1", Value: 1},
		{Key: "test-count-2", Value: 1},
		{Key: "test-count-3", Value: 1},
		{Key: "test-count-4", Value: 1},
		{Key: "test-count-5", Value: 1},
	},
}
var pc2 = payload.Payload{
	CountMetrics: []payload.CountMetric{
		{Key: "test-count-6", Value: 1},
		{Key: "test-count-1", Value: 1},
		{Key: "test-count-7", Value: 1},
		{Key: "test-count-2", Value: 1},
		{Key: "test-count-8", Value: 1},
	},
}
var pc3 = payload.Payload{
	CountMetrics: []payload.CountMetric{
		{Key: "test-count-7", Value: 1},
		{Key: "test-count-9", Value: 1},
		{Key: "test-count-8", Value: 1},
		{Key: "test-count-1", Value: 1},
		{Key: "test-count-10", Value: 1},
	},
}

func BenchmarkCounts(b *testing.B) {
	statsum, err := New(Config{JwtSecret: []byte("secret")})
	nilOrPanic(err, "Failed to create statsum")
	authorization := auth("p")
	p1 := jsonBytes(pc1)
	p2 := jsonBytes(pc2)
	p3 := jsonBytes(pc3)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p1))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", authorization)
			w := httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 1")

			r, _ = http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p2))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", authorization)
			w = httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 2")

			r, _ = http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p3))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", authorization)
			w = httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 3")
		}
	})
}

var pv1 = payload.Payload{
	ValueMetrics: []payload.ValueMetric{
		{Key: "test-count-1", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-2", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-3", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-4", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-5", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
	},
}
var pv2 = payload.Payload{
	ValueMetrics: []payload.ValueMetric{
		{Key: "test-count-6", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-1", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-7", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-2", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-8", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
	},
}
var pv3 = payload.Payload{
	ValueMetrics: []payload.ValueMetric{
		{Key: "test-count-7", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-9", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-8", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-1", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
		{Key: "test-count-10", Value: []float64{1, 2, 3, 4, 5, 6, 7, 8, 10}},
	},
}

func BenchmarkValues(b *testing.B) {
	statsum, err := New(Config{JwtSecret: []byte("secret")})
	nilOrPanic(err, "Failed to create statsum")
	authorization := auth("p")
	p1 := jsonBytes(pv1)
	p2 := jsonBytes(pv2)
	p3 := jsonBytes(pv3)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p1))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", authorization)
			w := httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 1")

			r, _ = http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p2))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", authorization)
			w = httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 2")

			r, _ = http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p3))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", authorization)
			w = httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 3")
		}
	})
}

func BenchmarkValuesMsgPack(b *testing.B) {
	statsum, err := New(Config{JwtSecret: []byte("secret")})
	nilOrPanic(err, "Failed to create statsum")
	authorization := auth("p")

	p1 := msgpBytes(pv1)
	p2 := msgpBytes(pv2)
	p3 := msgpBytes(pv3)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, _ := http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p1))
			r.Header.Set("Content-Type", "application/msgpack")
			r.Header.Set("Authorization", authorization)
			w := httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 1")

			r, _ = http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p2))
			r.Header.Set("Content-Type", "application/msgpack")
			r.Header.Set("Authorization", authorization)
			w = httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 2")

			r, _ = http.NewRequest("POST", "https://statsum.local/v1/project/p", bytes.NewReader(p3))
			r.Header.Set("Content-Type", "application/msgpack")
			r.Header.Set("Authorization", authorization)
			w = httptest.NewRecorder()
			statsum.handler(w, r)
			assert(w.Code == http.StatusOK, "Failed to send request 3")
		}
	})
}
