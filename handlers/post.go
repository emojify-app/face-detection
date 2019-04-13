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

	"github.com/emojify-app/face-detection/logging"
	"gocv.io/x/gocv"
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
}

// NewPost creates a new face processor handler
func NewPost(cl string) *Post {
	return &Post{cascadeLocation: cl}
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

	fp := NewFaceProcessor(p.cascadeLocation)
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

// FaceProcessor detects the position of a face from an input image
type FaceProcessor struct {
	faceclassifier *gocv.CascadeClassifier
}

// NewFaceProcessor loads the cascades and creates a new face processor
// to clear memory Close() must be called
func NewFaceProcessor(cl string) *FaceProcessor {
	// create clasifiers, to avoid leaking memory classifiers must be
	// closed after use
	classifier1 := gocv.NewCascadeClassifier()
	classifier1.Load(cl + "/haarcascade_frontalface_default.xml")

	return &FaceProcessor{
		faceclassifier: &classifier1,
	}
}

// Close frees allocated memory
func (fp *FaceProcessor) Close() {
	fp.faceclassifier.Close()
}

// DetectFaces detects faces in the image and returns an array of rectangle
func (fp *FaceProcessor) DetectFaces(file string) (faces []image.Rectangle, bounds image.Rectangle, err error) {
	img := gocv.IMRead(file, gocv.IMReadColor)
	if img.Empty() {
		return nil, image.Rectangle{}, fmt.Errorf("Unable to read image")
	}
	defer img.Close()

	// convert the image to greyscale for better processing
	greyImg := gocv.NewMat()
	defer greyImg.Close()

	gocv.CvtColor(img, &greyImg, gocv.ColorRGBToGray)
	gocv.EqualizeHist(greyImg, &greyImg)

	// define the bounds of the image
	bds := image.Rectangle{Min: image.Point{}, Max: image.Point{X: img.Cols(), Y: img.Rows()}}

	// detect faces
	tmpfaces := fp.faceclassifier.DetectMultiScaleWithParams(
		greyImg, 1.07, 6, 0, image.Point{X: 10, Y: 10}, image.Point{X: 500, Y: 500},
	)

	return tmpfaces, bds, nil
}
