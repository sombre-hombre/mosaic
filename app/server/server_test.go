package server

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

func Test_CreateMosaic(t *testing.T) {
	ts := httptest.NewServer(routes())
	defer ts.Close()

	input, err := os.Open("../../img/samples/Owl.jpg")
	require.NoError(t, err)
	defer input.Close()

	req, err := http.NewRequest(
		http.MethodPost,
		ts.URL+"/api/v1/libraries/abstract/mosaics",
		input,
	)
	req.Header.Add("Content-Type", "image/jpeg")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, resp.StatusCode, http.StatusOK)
	require.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")

	respData, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.DetectContentType(respData), "image/jpeg")
}
