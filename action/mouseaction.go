package action

import (
	con "autoclick/constant"
	"autoclick/controller"
	"autoclick/model"
	"autoclick/pkg/messagebus"
	"fmt"

	hook "github.com/robotn/gohook"
)

func init() {
	fmt.Println("mouseaction init")
	mob := &MouseEventObserver{}
	messagebus.RegisterObserver(con.MouseEventObserverName, mob)

	hob := &HookObserver{}
	messagebus.RegisterObserver(con.HookObserverName, hob)
}

type MouseEventObserver struct {
}

func (*MouseEventObserver) OnEvent(ev interface{}) {
	channel, ok := ev.(chan hook.Event)
	if !ok {
		return
	}
	go checkMouseEvent(channel)
}

type HookObserver struct {
	isStarted bool
}

func (ob *HookObserver) OnEvent(ev interface{}) {
	state, ok := ev.(string)

	if !ok {
		return
	}
	if state == "start" && !ob.isStarted {
		channel := hook.Start()
		ob.isStarted = true
		messagebus.SendMsg(con.MouseEventObserverName, channel)
	} else if state == "stop" && ob.isStarted {
		hook.End()
		ob.isStarted = false
	}
}

func checkMouseEvent(channel chan hook.Event) {
	state := "init"
	for ev := range channel {
		//fmt.Printf("%+v\n", ev)
		if ev.Kind == hook.MouseDrag && state == "init" {
			state = "mouseDown"
			axis := model.Axis{
				Left:   int(ev.X),
				Top:    int(ev.Y),
				Right:  0,
				Bottom: 0,
			}
			controller.OnMouseDown(axis)
		} else if ev.Kind == hook.MouseDrag && state == "mouseDown" {
			axis := model.Axis{
				Right:  int(ev.X),
				Bottom: int(ev.Y),
				Left:   0,
				Top:    0,
			}
			controller.OnMouseMove(axis)
		} else if ev.Kind == hook.MouseDown && state == "mouseDown" {
			state = "init"
			axis := model.Axis{
				Right:  int(ev.X),
				Bottom: int(ev.Y),
				Left:   0,
				Top:    0,
			}
			controller.OnMouseUp(axis)
		} else {
			state = "init"
			//controller.OnUnknownMouseState(ev.String())
		}
	}
}
