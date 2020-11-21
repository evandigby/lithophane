package main

import (
	"flag"
	"image/png"
	"os"

	"github.com/hschendel/stl"
)

var (
	imageFile        = flag.String("imageFile", "", "Image file to load (png)")
	outputFile       = flag.String("outputFile", "", "Output file (stl)")
	minimumThickness = flag.Float64("minthickmm", 0.8, "Minimum thickness in MM")
	maximumThickness = flag.Float64("maxthickmm", 3.4, "Maximum thickness")
	maxSizeMMFlag    = flag.Float64("maxmm", 100, "Max size in millimetres")
	borderWidth      = flag.Float64("borderwidthmm", 4, "Border width in MM")
	borderDepth      = flag.Float64("borderdepthmm", 0, "Border depth in MM (if zero means same as thickness)")
	pixelsPerMMFlag  = flag.Int("ppmm", 5, "Pixels per Millimetre")
)

func main() {
	flag.Parse()

	pngFile, err := os.Open(*imageFile)
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(pngFile)
	if err != nil {
		panic(err)
	}

	imgSize := img.Bounds()
	deltX := float32(imgSize.Dx())
	deltY := float32(imgSize.Dy())

	maxSizeMM := float32(*maxSizeMMFlag)
	pixelsPerMM := float32(*pixelsPerMMFlag)

	wh := maxSizeMM / deltX
	if deltY > deltX {
		wh = maxSizeMM / deltY
	}
	xyScale := wh * pixelsPerMM

	img = scaleImage(img, xyScale)

	bw := imgToGray(img)

	realWidth := (deltX * xyScale) / pixelsPerMM
	realHeight := (deltY * xyScale) / pixelsPerMM

	borderThick := *borderDepth
	if borderThick == 0 {
		borderThick = *maximumThickness
	}

	meshTriangles := imgToMesh(bw, float32(*minimumThickness), float32(*maximumThickness), realWidth, realHeight, float32(borderThick), float32(*borderWidth))

	lithophaneFace := stl.Solid{
		Name:      "Lithophane",
		IsAscii:   false,
		Triangles: meshTriangles,
	}

	lithophaneFace.RecalculateNormals()

	f, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lithophaneFace.WriteAll(f)
}
