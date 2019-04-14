package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/emojify-app/face-detection/detection"
	"github.com/emojify-app/face-detection/logging"
)

// Request is base64 encoded image

// Response for the function
type Response struct {
	Faces  []image.Rectangle
	Bounds image.Rectangle
}

// Post is a http handler which detects faces using the OpenCV library and GoCV
type Post struct {
	cascadeLocation string
	logger          logging.Logger
	scaleFactor     float64
	minNeighbors    int
}

// NewPost creates a new face processor handler with default parameters
func NewPost(cascadeLocation string) *Post {
	return &Post{
		cascadeLocation: cascadeLocation,
		scaleFactor:     1.05,
		minNeighbors:    8,
	}
}

// NewPostWithParams creates a new face processor handler
// Face detection parameters can be tuned by setting scaleFactor and minNeighbors
func NewPostWithParams(cascadeLocation string, scaleFactor float64, minNeighbors int) *Post {
	return &Post{
		cascadeLocation: cascadeLocation,
		scaleFactor:     scaleFactor,
		minNeighbors:    minNeighbors,
	}
}

// ServeHTTP handles the request
func (p *Post) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// OpenCV can panic, make sure we catch it
	defer func() {
		if r := recover(); r != nil {
			p.logger.Log().Error("Recovered in f", "trace", r)
		}
	}()

	var data []byte
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "No response body", http.StatusBadRequest)
		return
	}

	typ := http.DetectContentType(data)
	if typ != "image/jpeg" && typ != "image/png" {
		http.Error(
			rw,
			"Only jpeg or png images, either raw uncompressed bytes or base64 encoded are acceptable inputs, you uploaded: "+typ,
			http.StatusBadRequest,
		)
		return
	}

	// create a temporary file to be read by OpenCV
	tmpfile, err := ioutil.TempFile("/tmp", "image")
	if err != nil {
		p.logger.Log().Error("Unable to create temporary file", "error", err)
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up temporary file on exit

	// copy the body into a temporary file
	io.Copy(tmpfile, bytes.NewBuffer(data))

	fp := detection.New(p.cascadeLocation, p.scaleFactor, p.minNeighbors)
	defer fp.Close()

	faces, bounds, err := fp.DetectFaces(tmpfile.Name())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Response{
		Faces:  faces,
		Bounds: bounds,
	}

	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error encoding output: %s", err), http.StatusInternalServerError)
		return
	}

	// return the coordinates
	rw.Write(j)
}
