package safeclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultTimeout(t *testing.T) {
	t.Skip("Skipping: not yet implemented — remove skip when development is complete")
	c := New()
	assert.Equal(t, 10*time.Second, c.timeout)
}

func TestWithTimeout(t *testing.T) {
	t.Skip("Skipping: not yet implemented — remove skip when development is complete")
	c := New(WithTimeout(5 * time.Second))
	assert.Equal(t, 5*time.Second, c.timeout)
	assert.Equal(t, 5*time.Second, c.http.Timeout)
}

func TestGet_Success(t *testing.T) {
	t.Skip("Skipping: not yet implemented — remove skip when development is complete")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	}))
	defer ts.Close()

	c := New()
	body, err := c.Get(ts.URL)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(body))
}

func TestGet_404ReturnsError(t *testing.T) {
	t.Skip("Skipping: not yet implemented — remove skip when development is complete")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer ts.Close()

	c := New()
	_, err := c.Get(ts.URL)
	require.Error(t, err, "a 404 must be returned as an error by the safe client")
}

func TestGet_500ReturnsError(t *testing.T) {
	t.Skip("Skipping: not yet implemented — remove skip when development is complete")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New()
	_, err := c.Get(ts.URL)
	require.Error(t, err, "a 500 must be returned as an error by the safe client")
}
