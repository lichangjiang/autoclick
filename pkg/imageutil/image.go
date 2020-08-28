package imageutil

import (
	"autoclick/model"
	"fmt"
	"image"
	"time"

	//	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"github.com/vitali-fedulov/images"

	"github.com/go-vgo/robotgo"
)

func CaptureImage(left, top, right, bottom int) (image.Image, error) {

	min := image.Point{
		left,
		top,
	}

	max := image.Point{
		right,
		bottom,
	}

	bounds := image.Rectangle{
		min,
		max,
	}

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ImageSimiliar(a, b image.Image) bool {
	// Calculate hashes and image sizes.
	hashA, imgSizeA := images.Hash(a)
	hashB, imgSizeB := images.Hash(b)

	// Image comparison.
	if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
		return true
	} else {
		return false
	}
}

func CreateImageForEvents(events map[string]model.Event) (map[string]model.Event, error) {
	for _, event := range events {
		img, err := CaptureImage(event.Axis.Left,
			event.Axis.Top, event.Axis.Right, event.Axis.Bottom)

		if err != nil {
			return nil, err
		}

		event.Image = img
		events[event.Name] = event
	}

	return events, nil
}

func StartOneEventStreamCheck(events map[string]model.Event, name string) error {
	eventName := name
LOOP:
	event := events[eventName]
	//fmt.Printf("%s start event\n", event.Name)
	if event.Image == nil {
		return fmt.Errorf("can not find event:%s", event.Name)
	}
	currentImg, err := CaptureImage(event.Axis.Left,
		event.Axis.Top, event.Axis.Right, event.Axis.Bottom)

	if err != nil {
		return err
	}

	similar := ImageSimiliar(event.Image, currentImg)
	if similar && !event.Revert {
		fmt.Printf("%s event work\n", event.Name)
		x := (event.Axis.Left + event.Axis.Right) / 2
		y := (event.Axis.Top + event.Axis.Bottom) / 2
		ox, oy := robotgo.GetMousePos()
		robotgo.MoveMouse(x, y)
		robotgo.MouseClick()
		robotgo.MoveMouse(ox, oy)
		time.Sleep(100 * time.Millisecond)
		eventName = event.NextEvent
		if eventName != "" {
			goto LOOP
		}
	} else if !similar && event.Revert {
		fmt.Printf("%s event work\n", event.Name)
		x := (event.Axis.Left + event.Axis.Right) / 2
		y := (event.Axis.Top + event.Axis.Bottom) / 2
		robotgo.MoveMouse(x, y)
		robotgo.MouseClick()
		robotgo.MoveMouse(0, 0)
		time.Sleep(10 * time.Millisecond)
		eventName = event.NextEvent
		if eventName != "" {
			goto LOOP
		}
	} else {
		fmt.Printf("%s event not work\n", event.Name)
	}

	return nil
}

func StartImageCheck(events map[string]model.Event,
	startEvents map[string]bool) error {
	for k, _ := range startEvents {
		if err := StartOneEventStreamCheck(events, k); err != nil {
			fmt.Printf("check image error:%+v\n", err)
		}
	}
	return nil
}
