package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
	"path"

	"github.com/fogleman/gg"
	"github.com/golang/geo/r2"
	"github.com/golang/geo/s2"
)

func normalizePoints(points []r2.Point) []r2.Point {
	normalized := make([]r2.Point, 0)

	minX := points[0].X
	minY := points[0].Y
	maxX := points[0].X
	maxY := points[0].Y
	for _, point := range points {
		minX = math.Min(minX, point.X)
		minY = math.Min(minY, point.Y)
		maxX = math.Max(maxX, point.X)
		maxY = math.Max(maxY, point.Y)
	}

	for _, point := range points {
		p := r2.Point{
			X: (point.X - minX) / (maxX - minX),
			Y: (point.Y - minY) / (maxY - minY)}
		normalized = append(normalized, p)
	}

	return normalized
}

func main() {
	var imageSize int
	flag.IntVar(&imageSize, "imageSize", 1000, "image dimension in pixels")
	flag.Parse()
	gpxDirectory := flag.Arg(0)

	entries, err := ioutil.ReadDir(gpxDirectory)
	if err != nil {
		log.Fatal("Failed to read directory: ", err)
	}

	points := make([]r2.Point, 0)
	projection := s2.NewMercatorProjection(100)
	for _, file := range entries {
		if path.Ext(file.Name()) == ".gpx" {
			gpx, parseErr := ParseGpx(path.Join(gpxDirectory, file.Name()))
			if parseErr != nil {
				log.Fatal(file.Name(), parseErr)
			}

			for point := range gpx.GetPoints() {
				projected := projection.FromLatLng(point)
				points = append(points, projected)
			}
		}
	}

	normalized := normalizePoints(points)
	dc := gg.NewContext(imageSize, imageSize)
	dc.SetRGB(0, 0, 0)
	for _, point := range normalized {
		dc.DrawPoint(
			point.X*float64(imageSize),
			// FIXME why do I need to invert the Y axis to correct the map?
			float64(imageSize)-point.Y*float64(imageSize),
			1)
		dc.Fill()
	}
	dc.SavePNG("output.png")
}
