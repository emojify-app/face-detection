package detection

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// Face detects the position of a face from an input image
type Face struct {
	faceclassifier *gocv.CascadeClassifier
	scaleFactor    float64
	minNeighbors   int
}

// New loads the cascades and creates a new face processor
// to clear memory Close() must be called
func New(cl string, scaleFactor float64, minNeighbors int) *Face {
	// create clasifiers, to avoid leaking memory classifiers must be
	// closed after use
	c1 := gocv.NewCascadeClassifier()
	ok := c1.Load(cl + "/haarcascade_frontalface_default.xml")
	if !ok {
		panic("unable to load haar cascade face" + cl + "/haarcascade_frontalface_default.xml")
	}

	return &Face{
		faceclassifier: &c1,
		scaleFactor:    scaleFactor,
		minNeighbors:   minNeighbors,
	}
}

// Close frees allocated memory
func (fp *Face) Close() {
	fp.faceclassifier.Close()
}

// DetectFaces detects faces in the image and returns an array of rectangle
func (fp *Face) DetectFaces(file string) (faces []image.Rectangle, bounds image.Rectangle, err error) {
	return fp.detectFacesHAAR(file)
}

func (fp *Face) detectFacesHAAR(file string) (faces []image.Rectangle, bounds image.Rectangle, err error) {
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
		greyImg, fp.scaleFactor, fp.minNeighbors, 0, image.Point{X: 10, Y: 10}, image.Point{X: 500, Y: 500},
	)

	return tmpfaces, bds, nil
}
