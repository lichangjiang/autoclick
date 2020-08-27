package action

import (
	con "autoclick/constant"
	"autoclick/controller"
	"autoclick/model"
	"autoclick/pkg/messagebus"
	"fmt"
	"sync/atomic"

	hook "github.com/robotn/gohook"
)

var state atomic.Value

func init() {
	fmt.Println("mouseaction init")
	state.Store("init")
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
	s, ok := ev.(string)

	if !ok {
		return
	}
	if s == "start" && !ob.isStarted {
		channel := hook.Start()
		ob.isStarted = true
		messagebus.SendMsg(con.MouseEventObserverName, channel)
	} else if s == "stop" && ob.isStarted {
		hook.End()
		state.Store("init")
		ob.isStarted = false
	} else if s == "reset" {
		state.Store("init")
	}
}

func checkMouseEvent(channel chan hook.Event) {
	for ev := range channel {
		//fmt.Printf("%+v\n", ev)
		s := state.Load().(string)
		if ev.Kind == hook.MouseDrag && s == "init" {
			state.Store("mouseDown")
			axis := model.Axis{
				Left:   int(ev.X),
				Top:    int(ev.Y),
				Right:  0,
				Bottom: 0,
			}
			controller.OnMouseDown(axis)
		} else if ev.Kind == hook.MouseDrag && s == "mouseDown" {
			axis := model.Axis{
				Right:  int(ev.X),
				Bottom: int(ev.Y),
				Left:   0,
				Top:    0,
			}
			controller.OnMouseMove(axis)
		} else if ev.Kind == hook.MouseDown && s == "mouseDown" {
			state.Store("stop")
			axis := model.Axis{
				Right:  int(ev.X),
				Bottom: int(ev.Y),
				Left:   0,
				Top:    0,
			}
			controller.OnMouseUp(axis)
		}
	}
}
