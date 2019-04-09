package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/emojify-app/face-detection/logging"
	"gocv.io/x/gocv"
)

// Request is base64 encoded image

// Response for the function
type Response struct {
	Faces       []image.Rectangle
	Bounds      image.Rectangle
	ImageBase64 string
}

// Post is a http handler which detects faces using the OpenCV library and GoCV
type Post struct {
	cascadeLocation string
	logger          logging.Logger
}

// ServeHTTP handles the request
func (p *Post) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	var data []byte

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "No response body", http.StatusBadRequest)
		return
	}

	typ := http.DetectContentType(data)
	if typ != "image/jpeg" && typ != "image/png" {
		http.Error(rw, "Only jpeg or png images, either raw uncompressed bytes or base64 encoded are acceptable inputs, you uploaded: "+typ, http.StatusBadRequest)
	}

	tmpfile, err := ioutil.TempFile("/tmp", "image")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	io.Copy(tmpfile, bytes.NewBuffer(data))

	faceProcessor := NewFaceProcessor()
	faces, bounds := faceProcessor.DetectFaces(tmpfile.Name())

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

// BySize allows sorting images by size
type BySize []image.Rectangle

func (s BySize) Len() int {
	return len(s)
}
func (s BySize) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s BySize) Less(i, j int) bool {
	return s[i].Size().X > s[j].Size().X && s[i].Size().Y > s[j].Size().Y
}

var yellow = color.RGBA{255, 255, 0, 0}

// FaceProcessor detects the position of a face from an input image
type FaceProcessor struct {
	faceclassifier  *gocv.CascadeClassifier
	eyeclassifier   *gocv.CascadeClassifier
	glassclassifier *gocv.CascadeClassifier
}

// NewFaceProcessor creates a new face processor loading any dependent settings
func NewFaceProcessor() *FaceProcessor {
	// load classifier to recognize faces
	classifier1 := gocv.NewCascadeClassifier()
	classifier1.Load("./cascades/haarcascade_frontalface_default.xml")

	classifier2 := gocv.NewCascadeClassifier()
	classifier2.Load("./cascades/haarcascade_eye.xml")

	classifier3 := gocv.NewCascadeClassifier()
	classifier3.Load("./cascades/haarcascade_eye_tree_eyeglasses.xml")

	return &FaceProcessor{
		faceclassifier:  &classifier1,
		eyeclassifier:   &classifier2,
		glassclassifier: &classifier3,
	}
}

// DetectFaces detects faces in the image and returns an array of rectangle
func (fp *FaceProcessor) DetectFaces(file string) (faces []image.Rectangle, bounds image.Rectangle) {
	img := gocv.IMRead(file, gocv.IMReadColor)
	defer img.Close()

	bds := image.Rectangle{Min: image.Point{}, Max: image.Point{X: img.Cols(), Y: img.Rows()}}
	//gocv.CvtColor(img, img, gocv.ColorRGBToGray)
	//	gocv.Resize(img, img, image.Point{}, 0.6, 0.6, gocv.InterpolationArea)

	// detect faces
	tmpfaces := fp.faceclassifier.DetectMultiScaleWithParams(
		img, 1.07, 5, 0, image.Point{X: 10, Y: 10}, image.Point{X: 500, Y: 500},
	)

	fcs := make([]image.Rectangle, 0)

	if len(tmpfaces) > 0 {
		// draw a rectangle around each face on the original image
		for _, f := range tmpfaces {
			// detect eyes
			faceImage := img.Region(f)

			eyes := fp.eyeclassifier.DetectMultiScaleWithParams(
				faceImage, 1.01, 1, 0, image.Point{X: 0, Y: 0}, image.Point{X: 100, Y: 100},
			)

			if len(eyes) > 0 {
				fcs = append(fcs, f)
				continue
			}

			glasses := fp.glassclassifier.DetectMultiScaleWithParams(
				faceImage, 1.01, 1, 0, image.Point{X: 0, Y: 0}, image.Point{X: 100, Y: 100},
			)

			if len(glasses) > 0 {
				fcs = append(fcs, f)
				continue
			}
		}

		return fcs, bds
	}

	return nil, bds
}

// DrawFaces adds a rectangle to the given image with the face location
func (fp *FaceProcessor) DrawFaces(file string, faces []image.Rectangle) ([]byte, error) {
	if len(faces) == 0 {
		return ioutil.ReadFile(file)
	}

	img := gocv.IMRead(file, gocv.IMReadColor)
	defer img.Close()

	for _, r := range faces {
		gocv.Rectangle(&img, r, yellow, 1)
	}

	filename := fmt.Sprintf("/tmp/%d.jpg", time.Now().UnixNano())
	gocv.IMWrite(filename, img)
	defer os.Remove(filename) // clean up

	return ioutil.ReadFile(filename)
}
