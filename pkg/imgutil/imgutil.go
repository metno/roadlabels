package imgutil

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// Image with mean brightness < threshold is considered taken at night
const BRIGHTNESS_THRESHOLD = 0.38

type Image struct{ Img gocv.Mat }

func NewImage(imagePath string) (Image, error) {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return Image{}, fmt.Errorf("imgutil.NewImage(): Failed to load image %s", imagePath)
	}
	i := Image{Img: img}
	return i, nil
}

func NewImageFromBytes(bytebuf []byte) (Image, error) {
	img, err := gocv.IMDecode(bytebuf, gocv.IMReadColor)
	if err != nil {
		return Image{}, fmt.Errorf("imgutil.NewImage(): Failed to IMDecode image %s", err)
	}
	if img.Empty() {
		return Image{}, fmt.Errorf("imgutil.NewImage(): Failed to IMDecode image %s", err)
	}
	i := Image{Img: img}
	return i, nil
}

func (img Image) Close() {
	img.Img.Close()
}

// func (img Image) PutText(x int, y int, fontscale float64, thickness int, text string) {
func (img Image) PutText(text string) {
	color := color.RGBA{255, 0, 140, 0}

	var thickness, x, y int
	var fontscale float64
	//TODO: Calulate the following from image size some way or another:
	if img.Img.Rows() <= 256 { // Thumb
		thickness = 2
		fontscale = 0.6
		x = 10
		y = 25
	} else if img.Img.Rows() <= 400 {
		thickness = 1
		fontscale = 0.3
		x = 10
		y = 10
	} else if img.Img.Rows() <= 720 {
		thickness = 2
		fontscale = 0.8
		x = 10
		y = 50
	} else if img.Img.Rows() <= 2992 {
		thickness = 3
		fontscale = 3
		x = 10
		y = 100
	}

	pt := image.Point{x, y}
	gocv.PutText(&img.Img, text, pt, gocv.FontItalic, fontscale, color, thickness)
}

func (img Image) GetMeanBrightness() (float64, error) {

	// If we later need the hsv for other stuff, create it in NewImage()
	hsv := gocv.NewMat()
	defer hsv.Close()
	if img.Img.Empty() {

		return -1.0, fmt.Errorf("GetMeanBrightness: Bad matrix")
	}

	gocv.CvtColor(img.Img, &hsv, gocv.ColorBGRToHSV)
	return hsv.Mean().Val3 / 255.0, nil
}

// Wow. This is soo much easier in python-opencv
func (img Image) BlueMask(imagePath string) error {
	lowerBlue := gocv.NewMatFromScalar(gocv.NewScalar(102.0, 31.0, 160.0, 0.0), gocv.MatTypeCV8UC3)
	upperBlue := gocv.NewMatFromScalar(gocv.NewScalar(115.0, 255.0, 255.0, 0.0), gocv.MatTypeCV8UC3)
	hsv := gocv.NewMat()
	defer hsv.Close()
	if img.Img.Empty() {
		return fmt.Errorf("BlueMask: Bad matrix")

	}

	gocv.CvtColor(img.Img, &hsv, gocv.ColorBGRToHSV)

	channels, rows, cols := hsv.Channels(), hsv.Rows(), hsv.Cols()
	lowerChans := gocv.Split(lowerBlue)
	lowerMask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	lowerMaskChans := gocv.Split(lowerMask)
	// split HSV lower bounds into H, S, V channels
	upperChans := gocv.Split(upperBlue)
	upperMask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	upperMaskChans := gocv.Split(upperMask)

	// copy HSV values to upper and lower masks
	for c := 0; c < channels; c++ {
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				lowerMaskChans[c].SetUCharAt(i, j, lowerChans[c].GetUCharAt(0, 0))
				upperMaskChans[c].SetUCharAt(i, j, upperChans[c].GetUCharAt(0, 0))
			}
		}
	}

	gocv.Merge(lowerMaskChans, &lowerMask)
	gocv.Merge(upperMaskChans, &upperMask)
	// global mask

	mask := gocv.NewMat()
	defer mask.Close()
	gocv.InRange(hsv, lowerMask, upperMask, &mask)

	window := gocv.NewWindow("Hello")
	window.ResizeWindow(640, 480)
	window.IMShow(mask)
	window.WaitKey(-1)
	return nil
}

/*
func main() {
	path := "/lustre/storeB/project/metproduction/products/webcams/2018/07/04/136/136_20180704T1400Z.jpg"
	img, err := NewImage(path)
	if err != nil {
		fmt.Printf("Failed to load imgage %s: %v\n", path, err)
	}

	img.BlueMask(path)
}
*/
