package main

import (
	_ "autoclick/action"
	con "autoclick/constant"
	"autoclick/controller"
	"autoclick/pkg/messagebus"
	"autoclick/ui"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	myapp := app.New()
	appWin := myapp.NewWindow("auto click")

	msgLabel := widget.NewLabel("Welcome to autoclick")

	startBtn := widget.NewButton(con.StartBtnText, controller.OnStartBtnClick)

	addESBtn := widget.NewButton(con.AddEsBtnText, controller.OnAddEventStreamBtnClick)

	addEvBtn := widget.NewButton(con.AddEvBtnText, controller.OnAddEventBtnClick)

	resetBtn := widget.NewButton(con.ResetBtnText, controller.OnResetBtnClick)

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
	ui.SetBtn(con.StartBtnName, startBtn)
	ui.SetBtn(con.AddESBtnName, addESBtn)
	ui.SetBtn(con.AddEvBtnName, addEvBtn)
	addEvBtn.Disable()
	ui.SetBtn(con.ResetBtnName, resetBtn)

	appWin.SetOnClosed(func() {
		messagebus.CloseAll()
	})

	appWin.ShowAndRun()
}
