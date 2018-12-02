package main

import (
	"encoding/xml"
	"io/ioutil"
	"strconv"

	"github.com/golang/geo/s2"
)

type Waypoint struct {
	Latitude  string `xml:"lat,attr"`
	Longitude string `xml:"lon,attr"`
}

type Segment struct {
	Waypoints []Waypoint `xml:"trkpt"`
}

type Track struct {
	Segments []Segment `xml:"trkseg"`
}

type Gpx struct {
	Version string  `xml:"version,attr"`
	Tracks  []Track `xml:"trk"`
}

// GetPoints returns an iterator over all points in the GPX file.
func (g *Gpx) GetPoints() <-chan s2.LatLng {
	ch := make(chan s2.LatLng)
	go func() {
		for _, track := range g.Tracks {
			for _, segment := range track.Segments {
				for _, waypoint := range segment.Waypoints {
					lat, _ := strconv.ParseFloat(waypoint.Latitude, 64)
					lng, _ := strconv.ParseFloat(waypoint.Longitude, 64)
					latlng := s2.LatLngFromDegrees(lat, lng)
					ch <- latlng
				}
			}
		}
		close(ch)
	}()

	return ch
}

// ParseGpx returns track and waypoint information from a GPX file given by its
// filename.
func ParseGpx(filename string) (Gpx, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Gpx{}, err
	}

	parsed := Gpx{}
	err = xml.Unmarshal(data, &parsed)
	if err != nil {
		return Gpx{}, err
	}

	return parsed, nil
}
