package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	reader io.Reader
)

func setupTest() func(t *testing.T) {

	h := New()
	server = httptest.NewServer(h.endpoints())

	return func(t *testing.T) {
		server.Close()
	}
}

func TestHealth(t *testing.T) {
	tearDownTest := setupTest()
	defer tearDownTest(t)

	assert := assert.New(t)
	request, err := http.NewRequest("GET", fmt.Sprintf("%s", server.URL), nil)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(request)
	assert.Nil(err)
	assert.True(res.StatusCode == http.StatusOK)
}
