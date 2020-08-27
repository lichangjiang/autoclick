package controller

import (
	con "autoclick/constant"
	"autoclick/model"
	"autoclick/pkg/messagebus"
	"autoclick/ui"
	"fmt"

	"encoding/json"
)

func OnStartBtnClick() {
	btnText := ui.GetBtnText(con.StartBtnName)
	if btnText == con.StartBtnText {
		ui.ChangeBtnText(con.StartBtnName, "stop")
		ui.ShowMessage("auto click working...")
		ui.DisableAllOtherBtn(con.StartBtnName)

		msg := model.EventStreamMsg{
			Msg: "start",
		}
		messagebus.SendMsg(con.GlobalEventObserverName, msg)

		messagebus.SendMsg(con.UIStateObserverName,con.StartState)
	} else if btnText == "stop" {
		ui.EnableAllOtherBtn(con.AddEvBtnName)
		ui.ChangeBtnText(con.StartBtnName, con.StartBtnText)
		ui.ShowMessage("auto click stop")
		
		messagebus.SendMsg(con.UIStateObserverName,con.StopState)

		workMsg := model.WorkMsg {
			Msg : con.StopState,
		}
		messagebus.SendMsg(con.WorkObserverName,workMsg)
	}
}

func OnAddEventStreamBtnClick() {
	btnText := ui.GetBtnText(con.AddESBtnName)
	fmt.Printf("onAddEventStreamBtn Click text:%s\n", btnText)
	if btnText == con.AddEsBtnText {
		ui.ChangeBtnText(con.AddESBtnName, con.FinishBtnText)
		ui.EnableBtn(con.AddEvBtnName)
		ui.DisableAllOtherBtn(con.AddESBtnName, con.AddEvBtnName)
		ui.ShowMessage("click add event button,begin to record click event")
		msg := model.EventStreamMsg{
			Msg: con.NewStreamMsg,
		}
		messagebus.SendMsg(con.EventStreamObserverName, msg)
	} else if btnText == con.FinishBtnText {
		ui.DisableBtn(con.AddEvBtnName)
		ui.EnableAllOtherBtn(con.AddEvBtnName)
		ui.ChangeBtnText(con.AddESBtnName, con.AddEsBtnText)
		ui.ShowMessage(con.CommonText)
		msg := model.EventStreamMsg{
			Msg: con.EndStreamMsg,
		}
		messagebus.SendMsg(con.EventStreamObserverName, msg)
	}
}

func OnAddEventBtnClick() {
	btnText := ui.GetBtnText(con.AddEvBtnName)
	fmt.Printf("onAddEventBtn Click text:%s\n", btnText)
	if btnText == con.AddEvBtnText {
		ui.DisableBtn(con.AddESBtnName)
		ui.ChangeBtnText(con.AddEvBtnName, con.FinishBtnText)
		ui.ShowMessage("left: ,top: ,right ,bottom ")

		messagebus.SendMsg(con.HookObserverName, "start")
	} else if btnText == con.FinishBtnText {
		ui.EnableBtn(con.AddESBtnName)
		ui.ChangeBtnText(con.AddEvBtnName, con.AddEvBtnText)
		s, err := ui.GetShowMessage()
		if err != nil {
			fmt.Printf("GetShowMessage error:%+v\n", err)
			return
		}
		fmt.Println("json:" + s)
		axis := toAxis([]byte(s))
		if axis == nil {
			return
		}
		messagebus.SendMsg(con.HookObserverName, "stop")
		msg := model.EventStreamMsg{
			Msg:   con.AddEventMsg,
			Value: *axis,
		}
		messagebus.SendMsg(con.EventStreamObserverName, msg)

		ui.ShowMessage("add new event or finish")
	}
}

func OnResetBtnClick() {
	ui.ShowMessage("reset setting")
}

func OnMouseDown(axis model.Axis) {
	showAxis(axis)
}

func OnMouseMove(axis model.Axis) {
	axisB, err := ui.GetShowMessage()
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	oldAxis := toAxis([]byte(axisB))
	if oldAxis == nil {
		return
	}

	axis.Left = oldAxis.Left
	axis.Top = oldAxis.Top
	showAxis(axis)
}

func OnMouseUp(axis model.Axis) {
	OnMouseMove(axis)
}

func OnUnknownMouseState(msg string) {
	ui.ShowMessage("unknown mouse state:" + msg)
}

func showAxis(axis model.Axis) {
	str := axis2Str(axis)
	if str != "" {
		ui.ShowMessage(str)
	}
}

func axis2Str(axis model.Axis) string {
	b, err := json.Marshal(&axis)
	if err != nil {
		ui.ShowMessage(fmt.Sprintf("json error:%+v", err))
		return ""
	}
	return string(b)
}

func toAxis(b []byte) *model.Axis {
	oldAxis := &model.Axis{}
	err := json.Unmarshal([]byte(b), oldAxis)
	if err != nil {
		ui.ShowMessage(fmt.Sprintf("json error:%+v", err))
		return nil
	}
	return oldAxis
}
