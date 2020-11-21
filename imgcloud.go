package main

import (
	"image"
	"math"

	"github.com/hschendel/stl"
	"github.com/nfnt/resize"
)

func scaleImage(img image.Image, scale float32) image.Image {
	return resize.Resize(
		uint(float32(img.Bounds().Dx())*scale),
		uint(float32(img.Bounds().Dy())*scale),
		img,
		resize.Bilinear,
	)
}

func imgToGray(img image.Image) *image.Gray16 {
	bw := image.NewGray16(img.Bounds())

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			bw.Set(x, y, img.At(x, y))
		}
	}

	return bw
}

func imgToPoints(img *image.Gray16,
	minThick,
	maxThick,
	maxWidth,
	maxHeight float32,
) [][]stl.Vec3 {
	bounds := img.Bounds()
	min := bounds.Min

	// Min will usually but possibly not always be 0
	deltX := bounds.Dx()
	deltY := bounds.Dy()

	pixels := make([][]stl.Vec3, deltX)

	maxOffset := maxThick - minThick

	for x := 0; x < deltX; x++ {
		pixels[x] = make([]stl.Vec3, deltY)
		actualX := (float32(x) / float32(deltX)) * maxWidth
		for y := 0; y < deltY; y++ {
			pixel := img.Gray16At(min.X+x, min.Y+y)

			actualY := maxThick - ((float32(pixel.Y) / float32(math.MaxUint16)) * maxOffset)
			actualZ := (float32(y) / float32(deltY)) * maxHeight

			pixels[x][y] = vec(actualX, actualY, actualZ)
		}
	}

	return pixels
}

func makeWallPoints(row []stl.Vec3, toChange int, changer func(v float32) float32) []stl.Vec3 {
	newRow := make([]stl.Vec3, len(row))

	for i, v := range row {
		newRow[i] = stl.Vec3{v[0], v[1], v[2]}
		newRow[i][toChange] = changer(v[toChange])
	}

	return newRow
}

func makeWall(row1, row2 []stl.Vec3) []stl.Triangle {
	triangles := make([]stl.Triangle, 0)
	for i := 0; i < len(row1)-1; i++ {
		topLeft := row2[i]
		topRight := row1[i]
		bottomLeft := row2[i+1]
		bottomRight := row1[i+1]

		triangles = append(triangles, trianglesFromSquare(topLeft, bottomLeft, bottomRight, topRight)...)
	}
	return triangles
}

func borderForMesh(meshPoints [][]stl.Vec3, borderThick, borderWidth float32) [][]stl.Vec3 {
	zeroChanger := func(v float32) float32 { return 0 }
	borderThickChanger := func(v float32) float32 { return float32(borderThick) }
	wallDistChanger := func(sign float32) func(v float32) float32 {
		return func(v float32) float32 {
			return v + (borderWidth * sign)
		}
	}

	firstRow := meshPoints[0]
	lastRow := meshPoints[len(meshPoints)-1]

	firstRowWall := makeWallPoints(firstRow, 1, borderThickChanger)
	firstRowSide := makeWallPoints(firstRowWall, 0, wallDistChanger(-1))
	firstRowBack := makeWallPoints(firstRowSide, 1, zeroChanger)
	lastRowWall := makeWallPoints(lastRow, 1, borderThickChanger)
	lastRowSide := makeWallPoints(lastRowWall, 0, wallDistChanger(1))
	lastRowBack := makeWallPoints(lastRowSide, 1, zeroChanger)

	meshPoints = append([][]stl.Vec3{lastRowBack, firstRowBack, firstRowSide, firstRowWall}, meshPoints...)

	return append(meshPoints, lastRowWall, lastRowSide, lastRowBack)
}

func imgToMesh(img *image.Gray16,
	minThick,
	maxThick,
	maxWidth,
	maxHeight,
	borderThick,
	borderWidth float32,
) []stl.Triangle {
	points := imgToPoints(img, float32(*minimumThickness), float32(*maximumThickness), maxWidth, maxHeight)

	points = borderForMesh(points, borderThick, borderWidth)

	return meshPointsToMeshTriangles(points)
}
