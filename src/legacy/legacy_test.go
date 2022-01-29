package legacy

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirector(t *testing.T) {
	server := httptest.NewServer(NewRedirector(
		map[string]string{
			"/foo.html": "/bar.html",
		},
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, err := writer.Write([]byte("not filtered"))
			require.NoError(t, err)
		}),
	))
	defer server.Close()
	hc := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// consulting one of the paths marked for redirection results in a redirection
	r, err := hc.Get(server.URL + "/foo.html")
	require.NoError(t, err)
	assert.Equal(t, http.StatusMovedPermanently, r.StatusCode)
	assert.Equal(t, server.URL+"/bar.html", r.Header["Location"][0])

	// consulting any other path will forward the request to the underlying handler
	r, err = hc.Get(server.URL + "/index.html")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	body, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Equal(t, "not filtered", string(body))
}
