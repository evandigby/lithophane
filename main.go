package main

import (
	"os"

	"github.com/hschendel/stl"
)

func vec(x, y, z float32) stl.Vec3 {
	return stl.Vec3{x, y, z}
}

func main() {
	const (
		width  = 25.0
		height = 25.0
		depth  = 25.0
		yMin   = 0
		yMax   = depth
		xMin   = 0
		xMax   = width
		zMin   = 0
		zMax   = height
	)

	var (
		frontBottomLeft  = vec(xMin, yMin, zMin)
		frontBottomRight = vec(xMax, yMin, zMin)
		frontTopLeft     = vec(xMin, yMin, zMax)
		frontTopRight    = vec(xMax, yMin, zMax)

		backBottomLeft  = vec(xMin, yMax, zMin)
		backBottomRight = vec(xMax, yMax, zMin)
		backTopLeft     = vec(xMin, yMax, zMax)
		backTopRight    = vec(xMax, yMax, zMax)
	)

	solid := stl.Solid{
		Name:    "Cube",
		IsAscii: true,
		Triangles: []stl.Triangle{
			// Bottom
			{
				Normal:   vec(0, 0, -1),
				Vertices: [3]stl.Vec3{backBottomLeft, backBottomRight, frontBottomRight},
			},
			{
				Normal:   vec(0, 0, -1),
				Vertices: [3]stl.Vec3{frontBottomRight, frontBottomLeft, backBottomLeft},
			},
			// Front
			{
				Normal:   vec(0, 1, 0),
				Vertices: [3]stl.Vec3{frontBottomLeft, frontBottomRight, frontTopRight},
			},
			{
				Normal:   vec(0, 1, 0),
				Vertices: [3]stl.Vec3{frontTopRight, frontTopLeft, frontBottomLeft},
			},
			// Back
			{
				Normal:   vec(0, -1, 0),
				Vertices: [3]stl.Vec3{backBottomLeft, backTopLeft, backTopRight},
			},
			{
				Normal:   vec(0, -1, 0),
				Vertices: [3]stl.Vec3{backTopRight, backBottomRight, backBottomLeft},
			},
			// Left
			{
				Normal:   vec(-1, 0, 0),
				Vertices: [3]stl.Vec3{frontBottomLeft, frontTopLeft, backTopLeft},
			},
			{
				Normal:   vec(-1, 0, 0),
				Vertices: [3]stl.Vec3{backTopLeft, backBottomLeft, frontBottomLeft},
			},
			// Right
			{
				Normal:   vec(1, 0, 0),
				Vertices: [3]stl.Vec3{frontBottomRight, backBottomRight, backTopRight},
			},
			{
				Normal:   vec(1, 0, 0),
				Vertices: [3]stl.Vec3{backTopRight, frontTopRight, frontBottomRight},
			},
			// Top
			{
				Normal:   vec(0, 0, 1),
				Vertices: [3]stl.Vec3{frontTopLeft, frontTopRight, backTopRight},
			},
			{
				Normal:   vec(0, 0, 1),
				Vertices: [3]stl.Vec3{backTopRight, backTopLeft, frontTopLeft},
			},
		},
	}

	f, err := os.Create("test.stl")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	solid.WriteAll(f)
}
