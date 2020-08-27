package model

import "image"

type Event struct {
	Name      string      `json:"name"`
	Axis      Axis        `json:"axis"`
	Action    string      `json:"action"`
	NextEvent string      `json:"next"`
	Revert    bool        `json:"needDiff"`
	Image     image.Image `json:"-"`
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
	IsWrite bool
	EventMsg
}

type WorkMsg struct {
	Msg string
	EventMsg
}
