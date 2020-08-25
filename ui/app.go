package ui

import (
	"fmt"
	"sync/atomic"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type myApp struct {
	isStarted        atomic.Value
	isAddEventStream atomic.Value
	fyneApp          fyne.App
	window           fyne.Window
	messageLabel     *widget.Label
	btnMap           map[string]*widget.Button
}

var appIns = myApp{}

func init() {
	appIns.isStarted.Store(false)
	appIns.isAddEventStream.Store(false)
	appIns.btnMap = make(map[string]*widget.Button)
}

func SetApp(fyneApp fyne.App) {
	appIns.fyneApp = fyneApp
}

func SetWindow(win fyne.Window) {
	appIns.window = win
}

func SetMessageLabel(label *widget.Label) {
	appIns.messageLabel = label
}

func SetBtn(btnName string, btn *widget.Button) {
	appIns.btnMap[btnName] = btn
}

func DisableAllOtherBtn(keepBtnNames ...string) {
	fmt.Printf("keep btn:%+v\n", keepBtnNames)
	for name, btn := range appIns.btnMap {
		keep := false
		for _, keepName := range keepBtnNames {
			if name == keepName {
				keep = true
				break
			}
		}
		if !keep {
			fmt.Printf("disable btn:%s\n", name)
			btn.Disable()
		}
	}
}

func DisableBtn(btnName string) {
	if appIns.btnMap[btnName] != nil {
		appIns.btnMap[btnName].Disable()
	}
}

func EnableBtn(btnName string) {
	if appIns.btnMap[btnName] != nil {
		appIns.btnMap[btnName].Enable()
	}
}

func EnableAllBtn() {
	for _, btn := range appIns.btnMap {
		btn.Enable()
	}
}

func ChangeBtnText(btnName, text string) {
	if appIns.btnMap[btnName] != nil {
		appIns.btnMap[btnName].SetText(text)
	}
}

func GetBtnText(btnName string) string {
	if appIns.btnMap[btnName] != nil {
		return appIns.btnMap[btnName].Text
	}

	return ""
}

func GetBtnByName(btnName string) *widget.Button {
	return appIns.btnMap[btnName]
}

func ShowMessage(msg string) {
	if appIns.messageLabel != nil {
		appIns.messageLabel.SetText(msg)
	}
}
