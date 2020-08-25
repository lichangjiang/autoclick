package main

import (
	"autoclick/controller"
	"autoclick/ui"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	myapp := app.New()
	appWin := myapp.NewWindow("auto click")

	msgLabel := widget.NewLabel("Welcome to autoclick")

	startBtn := widget.NewButton(ui.StartBtnText,controller.OnStartBtnClick)

	addESBtn := widget.NewButton(ui.AddEsBtnText,controller.OnAddEventStreamBtnClick)

	addEvBtn := widget.NewButton(ui.AddEvBtnText,controller.OnAddEventBtnClick)

	resetBtn := widget.NewButton(ui.ResetBtnText,controller.OnResetBtnClick)

	appWin.SetContent(widget.NewVBox(
		msgLabel,
		widget.NewHBox(
			startBtn,
			addESBtn,
			addEvBtn,
			resetBtn,
		),
	))

	ui.SetApp(myapp)
	ui.SetWindow(appWin)
	ui.SetMessageLabel(msgLabel)
	ui.SetBtn(ui.StartBtnName,startBtn)
	ui.SetBtn(ui.AddESBtnName,addESBtn)
	ui.SetBtn(ui.AddEvBtnName,addEvBtn)
	ui.SetBtn(ui.ResetBtnName,resetBtn)

	appWin.ShowAndRun()
}
