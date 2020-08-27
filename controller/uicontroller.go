package controller

import (
	con "autoclick/constant"
	"autoclick/model"
	"autoclick/pkg/messagebus"
	"autoclick/ui"
	"fmt"
	"strings"

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

		messagebus.SendMsg(con.UIStateObserverName, con.StartState)
	} else if btnText == "stop" {
		ui.EnableAllOtherBtn(con.AddEvBtnName)
		ui.ChangeBtnText(con.StartBtnName, con.StartBtnText)
		ui.ShowMessage("auto click stop")

		messagebus.SendMsg(con.UIStateObserverName, con.StopState)

		workMsg := model.WorkMsg{
			Msg: con.StopState,
		}
		messagebus.SendMsg(con.WorkObserverName, workMsg)
	}
}

func OnAddEventStreamBtnClick() {
	btnText := ui.GetBtnText(con.AddESBtnName)
	fmt.Printf("onAddEventStreamBtn Click text:%s\n", btnText)
	if btnText == con.AddEsBtnText {
		ui.ChangeBtnText(con.AddESBtnName, con.FinishBtnText)
		ui.EnableBtn(con.AddEvBtnName)
		ui.DisableBtn(con.StartBtnName)
		ui.ChangeBtnText(con.ResetBtnName, con.ResetEventStreamBtnText)
		//ui.DisableAllOtherBtn(con.AddESBtnName, con.AddEvBtnName)
		ui.ShowMessage("click add event button,begin to record click event")
		msg := model.EventStreamMsg{
			Msg: con.NewStreamMsg,
		}
		messagebus.SendMsg(con.EventStreamObserverName, msg)
	} else if btnText == con.FinishBtnText {
		ui.DisableBtn(con.AddEvBtnName)
		ui.EnableAllOtherBtn(con.AddEvBtnName)
		ui.ChangeBtnText(con.AddESBtnName, con.AddEsBtnText)
		ui.ChangeBtnText(con.ResetBtnName, con.ResetBtnText)
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
		ui.ChangeBtnText(con.ResetBtnName, con.ResetEventBtxText)
		ui.ShowMessage("left: ,top: ,right ,bottom ")

		messagebus.SendMsg(con.HookObserverName, "start")
	} else if btnText == con.FinishBtnText {
		ui.EnableBtn(con.AddESBtnName)
		ui.ChangeBtnText(con.AddEvBtnName, con.AddEvBtnText)
		ui.ChangeBtnText(con.ResetBtnName, con.ResetEventStreamBtnText)
		s, err := ui.GetShowMessage()
		if err != nil {
			fmt.Printf("GetShowMessage error:%+v\n", err)
			return
		}
		fmt.Println("json:" + s)
		if !strings.HasPrefix(s, "left") {
			axis := toAxis([]byte(s))
			if axis != nil {
				messagebus.SendMsg(con.HookObserverName, "stop")
				msg := model.EventStreamMsg{
					Msg:   con.AddEventMsg,
					Value: *axis,
				}
				messagebus.SendMsg(con.EventStreamObserverName, msg)
			}
		}
		ui.ShowMessage("add new event or finish")
	}
}

func OnResetBtnClick() {
	btxText := ui.GetBtnText(con.ResetBtnName)

	if btxText == con.ResetBtnText {
		ui.ShowMessage("reset setting")
		messagebus.SendMsg(con.GlobalEventObserverName, model.EventStreamMsg{
			Msg: con.ResetEventStream,
		})
	} else if btxText == con.ResetEventStreamBtnText {
		ui.ShowMessage("reset current event stream")
		messagebus.SendMsg(con.EventStreamObserverName, model.EventStreamMsg{
			Msg: "reset",
		})
	} else if btxText == con.ResetEventBtxText {
		ui.ShowMessage("left: ,top: ,right ,bottom ")
		messagebus.SendMsg(con.HookObserverName, "reset")
	}
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
