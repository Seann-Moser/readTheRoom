package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"

	_ "embed"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed Poppins-Light.ttf
var fontData []byte

func cmToPixels(cm float64, dpi int) int {
	return int(cm * float64(dpi) / 2.54)
}

func drawDashedRectangle(dc *gg.Context, x1, y1, x2, y2 float64, dashLen, gapLen float64, lineWidth int, text string, face font.Face, drawBottom bool) {
	dc.Push()
	defer dc.Pop()

	// Set drawing parameters
	dc.SetLineWidth(float64(lineWidth))
	dc.SetColor(color.Black)
	dc.SetDash(dashLen, gapLen)

	// Draw rectangle sides using native dashed lines
	// Top side
	dc.DrawLine(x1, y1, x2, y1)
	// Right side
	dc.DrawLine(x2, y1, x2, y2)
	// Bottom side (if enabled)
	if drawBottom {
		dc.DrawLine(x2, y2, x1, y2)
	}
	// Left side
	dc.DrawLine(x1, y2, x1, y1)
	dc.Stroke()

	// Draw text
	dc.SetFontFace(face)
	cx := x1 + (x2-x1)/2
	cy := y1 + (y2-y1)/2
	dc.DrawStringAnchored(text, cx, cy, 0.5, 0.5)
}

func main() {
	dpi := 300
	width := cmToPixels(21.59, dpi)  // 8.5 inches in cm
	height := cmToPixels(27.94, dpi) // 11 inches in cm
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defaultSave := path.Join(wd, "ReadTheRoomImages")
	var c string
	fmt.Printf("Save Location(%s): ", defaultSave)
	fmt.Scanln(&c)
	if len(c) == 0 {
		c = defaultSave
	}

	for genFunc(c, width, height, dpi) {

	}
}

func genFunc(rootDir string, width int, height int, dpi int) bool {
	err := createDirIfNotExist(rootDir)
	if err != nil {
		return false
	}
	// Create image and context
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	dc := gg.NewContextForRGBA(img)

	// Load fonts
	//fontData, err := os.ReadFile("Poppins-Light.ttf")
	//if err != nil {
	//	panic(err)
	//}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}

	// Get central text
	var centralText string
	var sizeText string
	fmt.Print("Enter Word(or q to quit): ")
	fmt.Scanln(&centralText)
	if centralText == "q" {
		return false
	}
	fmt.Print("Font Size(default 120): ")

	fmt.Scanln(&sizeText)
	var size int
	size = 120
	if sizeText != "" {
		size, err = strconv.Atoi(sizeText)
		if err != nil {
			return false
		}
	}

	// Create rotated central text
	bigFace := truetype.NewFace(ttfFont, &truetype.Options{Size: 220})
	dc.SetFontFace(bigFace)
	textWidth, textHeight := dc.MeasureString(centralText)

	// Create separate image for text rotation
	textImg := image.NewRGBA(image.Rect(0, 0, int(textWidth), int(textHeight+90)))
	textDc := gg.NewContextForRGBA(textImg)
	textDc.SetFontFace(bigFace)
	textDc.SetColor(color.Black)
	textDc.DrawString(centralText, 0, textHeight)

	// Rotate text image
	rotated := gg.NewContext(textImg.Bounds().Dy(), textImg.Bounds().Dx())
	rotated.Rotate(gg.Radians(90))
	rotated.DrawImage(textImg, 0, -textImg.Bounds().Dy())

	// Calculate position and draw rotated text
	xPos := (width - rotated.Width()) / 2
	yPos := (height - rotated.Height()) / 2
	dc.DrawImage(rotated.Image(), xPos, yPos)

	// Draw dashed boxes
	rows := 12
	marginCm := 1.5
	widthCm := 21.59
	heightCm := 27.94
	leftXStartCm := marginCm
	rightXStartCm := widthCm - 1.5 - marginCm - 5
	boxHeightCm := (heightCm - 2*marginCm) / float64(rows)

	smallFace := truetype.NewFace(ttfFont, &truetype.Options{Size: float64(size)})

	for i := 0; i < rows; i++ {
		yTopCm := marginCm + float64(i)*boxHeightCm
		yBottomCm := yTopCm + boxHeightCm

		// Left box
		drawDashedRectangle(dc,
			float64(cmToPixels(leftXStartCm, dpi)),
			float64(cmToPixels(yTopCm, dpi)),
			float64(cmToPixels(leftXStartCm+6.55, dpi)),
			float64(cmToPixels(yBottomCm, dpi)),
			70, 40, 8, // dash/gap lengths in pixels
			centralText, smallFace, i == rows-1)

		// Right box
		drawDashedRectangle(dc,
			float64(cmToPixels(rightXStartCm, dpi)),
			float64(cmToPixels(yTopCm, dpi)),
			float64(cmToPixels(rightXStartCm+6.55, dpi)),
			float64(cmToPixels(yBottomCm, dpi)),
			70, 40, 8, // dash/gap lengths in pixels
			centralText, smallFace, i == rows-1)
	}

	// Save image
	fileName := path.Join(rootDir, strings.ReplaceAll(centralText, " ", "")+"_image.png")
	f, err := os.Create(fileName)
	if err != nil {
		return false
	}
	defer f.Close()
	png.Encode(f, img)
	fmt.Println("Saved", fileName)
	return true
}

func createDirIfNotExist(path string) error {
	// Check if the directory exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		fmt.Println("Directory created:", path)
	} else if err != nil {
		return fmt.Errorf("error checking directory: %w", err)
	} else {
		// Directory exists
		fmt.Println("Directory already exists:", path)
	}
	return nil
}
