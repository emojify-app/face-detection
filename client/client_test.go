package client

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/emojify-app/face-detection/handlers"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*httptest.Server, *HTTPClient) {
	s := httptest.NewServer(handlers.NewPost("../cascades"))

	return s, NewClient(s.URL)
}

func TestClientWorks(t *testing.T) {
	s, c := setup(t)
	defer s.Close()

	f, err := os.Open("../test_fixtures/group.jpg")
	assert.Nil(t, err)

	r, err := c.DetectFaces(f)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, 14, len(r.Faces))
}
