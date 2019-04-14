package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*httptest.ResponseRecorder, *http.Request, *Post) {
	return httptest.NewRecorder(), httptest.NewRequest("POST", "/", nil), NewPost("../cascades")
}

func TestReturnsBadRequestWhenBodyNotImage(t *testing.T) {
	input := "../test_fixtures/file.txt"
	rr, r, h := setup(t)
	data, err := ioutil.ReadFile(input)

	assert.Nil(t, err)

	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	h.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDetectsFacesInImageAndReturnsJSON(t *testing.T) {
	input := "../test_fixtures/group.jpg"
	rr, r, h := setup(t)
	data, err := ioutil.ReadFile(input)

	assert.Nil(t, err)

	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	j := &Response{}
	h.ServeHTTP(rr, r)
	json.Unmarshal(rr.Body.Bytes(), j)

	//fmt.Printf("%#v/n", j)
	assert.Equal(t, 14, len(j.Faces))
}
