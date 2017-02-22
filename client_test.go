package statsum

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taskcluster/statsum/server"
)

func TestClientAgainstServer(t *testing.T) {
	statsum, err := server.New(server.Config{
		Port:      "50321",
		JwtSecret: []byte("test-secret"),
	})
	require.NoError(t, err)
	go statsum.Start()

	client := New("test-project", StaticConfigurer(
		"http://localhost:50321",
		[]byte("test-secret"),
	), Options{
		MaxDataPoints: 10,
		MaxDelay:      2 * time.Second,
		MinDelay:      1 * time.Second,
		OnError:       func(err error) { panic(fmt.Sprintf("ASYNC ERROR: %s", err)) },
	})

	t.Run("Count", func(t *testing.T) {
		client.Count("test", 1)
		client.cache.mData.Lock()
		assert.Contains(t, client.cache.counters, "test")
		assert.EqualValues(t, client.cache.counters["test"], 1)
		assert.EqualValues(t, client.cache.dataPoints, 1)
		client.cache.mData.Unlock()

		client.Count("test", 2)
		client.cache.mData.Lock()
		assert.Contains(t, client.cache.counters, "test")
		assert.EqualValues(t, client.cache.counters["test"], 3)
		assert.EqualValues(t, client.cache.dataPoints, 2)
		client.cache.mData.Unlock()

		require.NoError(t, client.Flush())
		client.cache.mData.Lock()
		assert.NotContains(t, client.cache.counters, "test")
		assert.EqualValues(t, client.cache.dataPoints, 0)
		client.cache.mData.Unlock()
	})

	t.Run("Measure", func(t *testing.T) {
		client.Measure("test", 1.3)
		client.cache.mData.Lock()
		assert.Contains(t, client.cache.measures, "test")
		assert.EqualValues(t, client.cache.measures["test"], []float64{1.3})
		assert.EqualValues(t, client.cache.dataPoints, 1)
		client.cache.mData.Unlock()

		client.Measure("test", 2, 3)
		client.cache.mData.Lock()
		assert.Contains(t, client.cache.measures, "test")
		assert.EqualValues(t, client.cache.measures["test"], []float64{1.3, 2, 3})
		assert.EqualValues(t, client.cache.dataPoints, 3)
		client.cache.mData.Unlock()

		require.NoError(t, client.Flush())
		client.cache.mData.Lock()
		assert.NotContains(t, client.cache.measures, "test")
		assert.EqualValues(t, client.cache.dataPoints, 0)
		client.cache.mData.Unlock()
	})

	t.Run("Time", func(t *testing.T) {
		client.Time("test", func() {
			time.Sleep(500 * time.Millisecond)
		})
		client.cache.mData.Lock()
		assert.Contains(t, client.cache.measures, "test")
		assert.InDelta(t, 500, client.cache.measures["test"][0], 100)
		assert.EqualValues(t, client.cache.dataPoints, 1)
		client.cache.mData.Unlock()

		require.NoError(t, client.Flush())
		client.cache.mData.Lock()
		assert.NotContains(t, client.cache.measures, "test")
		assert.EqualValues(t, client.cache.dataPoints, 0)
		client.cache.mData.Unlock()
	})

	t.Run("WithPrefix", func(t *testing.T) {
		c := client.WithPrefix("this")
		c.Count("that", 42)
		c.Measure("that", 1.4, 5.4)
		c.cache.mData.Lock()
		assert.Contains(t, c.cache.counters, "this.that")
		assert.EqualValues(t, c.cache.counters["this.that"], 42)
		assert.Contains(t, c.cache.measures, "this.that")
		assert.EqualValues(t, c.cache.measures["this.that"], []float64{1.4, 5.4})
		assert.EqualValues(t, c.cache.dataPoints, 3)
		c.cache.mData.Unlock()

		c2 := c.WithPrefix("level2")
		c2.Count("that", 42)
		c2.Measure("that", 1.4)
		c2.cache.mData.Lock()
		assert.Contains(t, c2.cache.counters, "this.level2.that")
		assert.EqualValues(t, c2.cache.counters["this.level2.that"], 42)
		assert.Contains(t, c2.cache.measures, "this.level2.that")
		assert.EqualValues(t, c2.cache.measures["this.level2.that"], []float64{1.4})
		assert.EqualValues(t, c2.cache.dataPoints, 5)
		c2.cache.mData.Unlock()

		require.NoError(t, c.Flush())
		c.cache.mData.Lock()
		assert.NotContains(t, c.cache.measures, "test.that")
		assert.EqualValues(t, c.cache.dataPoints, 0)
		c.cache.mData.Unlock()
	})

	t.Run("MaxDataPoints", func(t *testing.T) {
		client.Measure("that", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
		client.cache.mData.Lock()
		assert.EqualValues(t, client.cache.dataPoints, 0)
		client.cache.mData.Unlock()

		require.NoError(t, client.Flush())
		client.cache.mData.Lock()
		assert.NotContains(t, client.cache.measures, "test")
		assert.EqualValues(t, client.cache.dataPoints, 0)
		client.cache.mData.Unlock()
	})
}

func TestAutomaticClientSubmission(t *testing.T) {
	msg := make(chan struct{})
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg <- struct{}{}
	}))
	defer s.Close()

	client := New("test-project", StaticConfigurer(
		s.URL,
		[]byte("test-secret"),
	), Options{
		MaxDataPoints: 3,
		MaxDelay:      4 * time.Second,
		MinDelay:      3 * time.Second,
		OnError:       func(err error) { panic(fmt.Sprintf("ASYNC ERROR: %s", err)) },
	})

	// Write 3 data points now we expect a flush, if not we hang because of a bug
	client.Count("test", 1)
	client.Count("test", 2)
	client.Count("test", 3)
	expected := time.Now()
	<-msg
	assert.WithinDuration(t, expected, time.Now(), 200*time.Millisecond)

	// Write 3 data points now we expect a flush, if not we hang because of a bug
	client.Measure("test", 1)
	client.Measure("test", 2.4, 6.7, 7.8)
	expected = time.Now()
	<-msg
	assert.WithinDuration(t, expected, time.Now(), 200*time.Millisecond)

	client.Measure("test2", 2323)
	expected = time.Now().Add(3500 * time.Millisecond)
	<-msg
	assert.WithinDuration(t, expected, time.Now(), 1100*time.Millisecond)

	// check that we didn't automagically send without any metrics
	time.Sleep(4 * time.Second)
	select {
	case <-msg:
		panic("Didn't expect automagic submissions")
	default:
	}
}
