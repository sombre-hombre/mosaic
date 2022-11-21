package server

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Ping(t *testing.T) {
    ts := httptest.NewServer(routes())
    defer ts.Close()

    resp, err := http.Get(ts.URL + "/internal/ping")
    require.NoError(t, err)
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    require.NoError(t, err)

    require.Equal(t, "pong", string(body))
    require.Equal(t, http.StatusOK, resp.StatusCode)
}
