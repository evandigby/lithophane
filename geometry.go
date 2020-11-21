package main

import "github.com/hschendel/stl"

func vec(x, y, z float32) stl.Vec3 {
	return stl.Vec3{x, y, z}
}

func trianglesFromSquare(first, second, third, fourth stl.Vec3) []stl.Triangle {
	return []stl.Triangle{
		{
			Vertices: [3]stl.Vec3{first, second, third},
		},
		{
			Vertices: [3]stl.Vec3{third, fourth, first},
		},
	}
}

func meshPointsToMeshTriangles(meshPoints [][]stl.Vec3) []stl.Triangle {
	triangles := make([]stl.Triangle, 0)

	for x := 0; x < len(meshPoints)-1; x++ {
		for y := 0; y < len(meshPoints[x])-1; y++ {
			topLeft := meshPoints[x][y]
			topRight := meshPoints[x+1][y]
			bottomLeft := meshPoints[x][y+1]
			bottomRight := meshPoints[x+1][y+1]

			triangles = append(triangles, trianglesFromSquare(topLeft, bottomLeft, bottomRight, topRight)...)
		}
	}

	return triangles
}

func expandPoints(points ...[]stl.Vec3) []stl.Vec3 {
	expanded := make([]stl.Vec3, 0)

	for _, pointSlice := range points {
		for _, point := range pointSlice {
			expanded = append(expanded, point)
		}
	}

	return expanded
}

func expandVerts(tris ...[]stl.Triangle) []stl.Triangle {
	expanded := make([]stl.Triangle, 0)

	for _, triSlice := range tris {
		for _, tri := range triSlice {
			expanded = append(expanded, tri)
		}
	}

	return expanded
}

func cube(width, height, depth float32) stl.Solid {
	var (
		yMin             float32  = 0
		yMax             float32  = depth
		xMin             float32  = 0
		xMax             float32  = width
		zMin             float32  = 0
		zMax             float32  = height
		frontBottomLeft  stl.Vec3 = vec(xMin, yMin, zMin)
		frontBottomRight stl.Vec3 = vec(xMax, yMin, zMin)
		frontTopLeft     stl.Vec3 = vec(xMin, yMin, zMax)
		frontTopRight    stl.Vec3 = vec(xMax, yMin, zMax)
		backBottomLeft   stl.Vec3 = vec(xMin, yMax, zMin)
		backBottomRight  stl.Vec3 = vec(xMax, yMax, zMin)
		backTopLeft      stl.Vec3 = vec(xMin, yMax, zMax)
		backTopRight     stl.Vec3 = vec(xMax, yMax, zMax)
	)

	return stl.Solid{
		Name:    "Cube",
		IsAscii: true,
		Triangles: expandVerts(
			// Bottom
			trianglesFromSquare(backBottomLeft, backBottomRight, frontBottomRight, frontBottomLeft),
			// Front
			trianglesFromSquare(frontBottomLeft, frontBottomRight, frontTopRight, frontTopLeft),
			// Back
			trianglesFromSquare(backBottomLeft, backTopLeft, backTopRight, backBottomRight),
			// Left
			trianglesFromSquare(frontBottomLeft, frontTopLeft, backTopLeft, backBottomLeft),
			// Right
			trianglesFromSquare(frontBottomRight, backBottomRight, backTopRight, frontTopRight),
			// Top
			trianglesFromSquare(frontTopLeft, frontTopRight, backTopRight, backTopLeft),
		),
	}
}
