package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func setup(t *testing.T) (*httptest.ResponseRecorder, *http.Request, *Post) {
	return httptest.NewRecorder(), httptest.NewRequest("POST", "/", nil), NewPost("../cascades")
}

func TestDetectsFacesInImageAndReturnsJSON(t *testing.T) {
	input := "../test_fixtures/group.jpg"
	rr, r, h := setup(t)
	is := is.New(t)
	data, err := ioutil.ReadFile(input)
	is.NoErr(err) // Error should be nil

	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	j := &Response{}
	h.ServeHTTP(rr, r)
	json.Unmarshal(rr.Body.Bytes(), j)

	fmt.Printf("%#v/n", j)
	is.Equal(true, len(j.Faces) == 14)
}
