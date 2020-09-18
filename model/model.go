package model

import "image"

type Event struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Axis
	Image image.Image
}

type Axis struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

type StartMsg struct {
	EventTimeInterval      int       `json:"eventTimeInterval"`
	EventGroupTimeInterval int       `json:"eventGroupTimeInterval"`
	EventGroups            [][]Event `json:"eventGroups"`
}
