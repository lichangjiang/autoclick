package model

type ScreenInfo struct {
	Num int `json:"number"`
	Width int `json:"width"`
	Height int `json:"height"`
	Axis
}