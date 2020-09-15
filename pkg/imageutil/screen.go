package imageutil

import (
	"autoclick/model"

	"github.com/kbinani/screenshot"
)

func GetScreenNum() int {
	return screenshot.NumActiveDisplays()
}

func GetScreenInfo(num int) model.ScreenInfo {
	bounds := screenshot.GetDisplayBounds(num)
	return model.ScreenInfo{
		Num:    num,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Axis: model.Axis{
			Left:   bounds.Min.X,
			Top:    bounds.Min.Y,
			Right:  bounds.Max.X,
			Bottom: bounds.Max.Y,
		},
	}
}

func GetAllScreenInfo() []model.ScreenInfo {
	displayNum := GetScreenNum()
	screens := []model.ScreenInfo{}
	for i := 0; i < displayNum; i++ {
		screens = append(screens, GetScreenInfo(i))
	}
	return screens
}
