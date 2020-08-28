package model

import "image"

type Event struct {
	Name      string      `json:"name"`
	Axis      Axis        `json:"axis"`
	Action    string      `json:"action"`
	NextEvent string      `json:"next"`
	Revert    bool        `json:"needDiff"`
	Image     image.Image `json:"-"`
	ImageFile string      `json:"imageFile"`
}

type Axis struct {
	Left   int `json:left`
	Top    int `json:top`
	Right  int `json:right`
	Bottom int `json:bottom`
}

type EventGroup struct {
	StartEvents []string `json:startEvents`
	EndEvents   []string `json:endEvents`
	Events      []Event  `json:events`
}

type EventStream struct {
	Events []*Event
}

type EventStreamMsg struct {
	Msg   string
	Value interface{}
}

type EventMsg struct {
	StartEventNames []string
	EventMap        map[string]*Event
}

type JsonMsg struct {
	IsWrite    bool
	IsReadJson bool
	NeedCopy   bool
	EventMsg
	IsDir         bool
	IsImage       bool
	Image         image.Image
	ImageFileName string
	IsDeleteImage bool
}

type WorkMsg struct {
	Msg string
	EventMsg
}
