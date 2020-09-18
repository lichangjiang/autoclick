package imageutil

import (
	"autoclick/model"
	"fmt"
	"image"
	"sync"
	"time"

	//	"github.com/go-vgo/robotgo"
	"github.com/Nr90/imgsim"
	"github.com/kbinani/screenshot"
	"github.com/sirupsen/logrus"

	"github.com/go-vgo/robotgo"
	"github.com/google/uuid"

	multierror "github.com/hashicorp/go-multierror"
)

var ScreenshotMutex sync.Mutex

func CaptureScreen(num int) (image.Image, error) {
	ScreenshotMutex.Lock()
	defer ScreenshotMutex.Unlock()
	return screenshot.CaptureDisplay(num)
}

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

	ScreenshotMutex.Lock()
	defer ScreenshotMutex.Unlock()
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ImageSimiliar(a, b image.Image) bool {
	// Calculate hashes and image sizes.

	hashA := imgsim.AverageHash(a)
	hashB := imgsim.AverageHash(b)

	// Image comparison.
	if imgsim.Distance(hashA, hashB) <= 8 {
		return true
	} else {
		return false
	}
}

func CreateImageForEvents(events map[string]*model.Event) (map[string]*model.Event, error) {
	for _, event := range events {
		if event.Image == nil {
			img, err := CaptureImage(event.Axis.Left,
				event.Axis.Top, event.Axis.Right, event.Axis.Bottom)
			if err != nil {
				return nil, err
			}
			event.Image = img
		}

		events[event.Name] = event
	}

	return events, nil
}

func StartOneEventStreamCheck(eventGroup []model.Event, eventTimeInterval int) error {

	for _, event := range eventGroup {
		if event.Image == nil {
			return fmt.Errorf("event:%s has not image object", event.Name)
		}
		currentImg, err := CaptureImage(event.Axis.Left,
			event.Axis.Top, event.Axis.Right, event.Axis.Bottom)

		if err != nil {
			return err
		}

		similar := ImageSimiliar(event.Image, currentImg)
		if similar {
			logrus.WithFields(logrus.Fields{
				"name":   event.Name,
				"left":   event.Axis.Left,
				"top":    event.Axis.Top,
				"right":  event.Axis.Right,
				"bottom": event.Axis.Bottom,
			}).Info("event work")

			x := (event.Axis.Left + event.Axis.Right) / 2
			y := (event.Axis.Top + event.Axis.Bottom) / 2
			ox, oy := robotgo.GetMousePos()
			robotgo.MoveMouse(x, y)

			time.Sleep(time.Duration(eventTimeInterval) * time.Second)
			robotgo.MouseClick()
			robotgo.MoveMouse(ox, oy)
		} else {
			logrus.WithFields(logrus.Fields{
				"name":   event.Name,
				"left":   event.Axis.Left,
				"top":    event.Axis.Top,
				"right":  event.Axis.Right,
				"bottom": event.Axis.Bottom,
			}).Info("event not work")
			break
		}
	}

	return nil
}

func StartImageCheck(eventGroups [][]model.Event,
	eventGroupTimeInterval, eventTimeInterval int) error {
	//加入随机因素
	eventGroupMap := map[string][]model.Event{}
	for _, eventGroup := range eventGroups {
		key := uuid.New().String()
		eventGroupMap[key] = eventGroup
	}
	var errResult error
	for _, eventGroup := range eventGroupMap {
		if err := StartOneEventStreamCheck(eventGroup, eventTimeInterval); err != nil {
			multierror.Append(errResult, err)
		}
		time.Sleep(time.Duration(eventGroupTimeInterval) * time.Second)
	}
	return errResult
}
