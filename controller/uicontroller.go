package controller

import "autoclick/ui"

func OnStartBtnClick() {
	btnText := ui.GetBtnText(ui.StartBtnName) 
	if btnText == ui.StartBtnText {
		ui.ChangeBtnText(ui.StartBtnName,"stop")
		ui.ShowMessage("auto click working...")
		ui.DisableAllOtherBtn(ui.StartBtnName)
	} else if btnText == "stop" {
		ui.EnableAllBtn()
		ui.ChangeBtnText(ui.StartBtnName,ui.StartBtnText)
		ui.ShowMessage("auto click stop")
	}	
}

func OnAddEventStreamBtnClick() {
	btnText := ui.GetBtnText(ui.AddESBtnName)
	if btnText == ui.AddEsBtnText {
		ui.ChangeBtnText(ui.AddESBtnName, ui.FinishBtnText)
		ui.DisableAllOtherBtn(ui.AddESBtnName, ui.AddEvBtnName)
		ui.ShowMessage("click add event button,begin to record click event")
	} else if btnText == ui.FinishBtnText {
		ui.EnableAllBtn()
		ui.ChangeBtnText(ui.AddESBtnName, ui.AddEsBtnText)
		ui.ShowMessage(ui.CommonText)
	}
}

func OnAddEventBtnClick() {
	btnText := ui.GetBtnText(ui.AddEvBtnName)
	if btnText == ui.AddEvBtnText {
		ui.DisableBtn(ui.AddESBtnName)
		ui.ChangeBtnText(ui.AddEvBtnName,ui.FinishBtnText)
		ui.ShowMessage("left: ,top: ,right ,bottom ")
	} else if btnText == ui.FinishBtnText {
		ui.EnableBtn(ui.AddESBtnName)
		ui.ChangeBtnText(ui.AddEvBtnName,ui.AddEvBtnText)
		ui.ShowMessage("add new event or finish")		
	}
}

func OnResetBtnClick() {
	ui.ShowMessage("reset setting")
}
