package action

import (
	"autoclick/constant"
	"autoclick/model"
	"autoclick/pkg/messagebus"
	"fmt"
	"strconv"
	"time"
)

type EventStreamObserver struct {
	currentEventStream *model.EventStream
}

func init() {
	fmt.Println("EventStreamObserver init")
	ob := &EventStreamObserver{}
	messagebus.RegisterObserver(constant.EventStreamObserverName, ob)

	gob := &globalEventObserver{
		[]string{},
		map[string]*model.Event{},
	}
	messagebus.RegisterObserver(constant.GlobalEventObserverName, gob)
}

func (ob *EventStreamObserver) OnEvent(ev interface{}) {
	msg, ok := ev.(model.EventStreamMsg)

	if !ok {
		fmt.Println("EventStreamObserver Event fail cast to EventStreamMsg")
		return
	}

	if msg.Msg == constant.AddEventMsg {
		fmt.Printf("EventStreamObserver receive AddEventMsg:%+v\n", &msg)
		axis, ok := msg.Value.(model.Axis)
		if !ok {
			fmt.Println("AddEventMsg Value fail cast to Axis")
			return
		}

		if ob.currentEventStream == nil {
			fmt.Println("CurrentEventStream is nil")
			return
		}
		n := time.Now().Nanosecond()
		event := &model.Event{
			Name:      "event" + strconv.Itoa(n),
			Axis:      axis,
			Action:    "click",
			NextEvent: "",
			Revert:    false,
		}

		eventLen := len(ob.currentEventStream.Events)
		if eventLen > 0 {
			pre := ob.currentEventStream.Events[eventLen-1]
			pre.NextEvent = event.Name
		}
		ob.currentEventStream.Events = append(ob.currentEventStream.Events, event)
	} else if msg.Msg == constant.NewStreamMsg || msg.Msg == "reset" {
		fmt.Printf("EventStreamObserver %s\n", msg.Msg)
		cs := &model.EventStream{
			Events: []*model.Event{},
		}
		ob.currentEventStream = cs
	} else if msg.Msg == constant.EndStreamMsg {
		cs := ob.currentEventStream
		ob.currentEventStream = nil
		if cs != nil {
			newMsg := model.EventStreamMsg{
				Msg:   constant.AddEventStream,
				Value: cs,
			}
			messagebus.SendMsg(constant.GlobalEventObserverName, newMsg)
		}
	}
}

type globalEventObserver struct {
	startEventNames []string
	eventMap        map[string]*model.Event
}

func (ob *globalEventObserver) OnEvent(ev interface{}) {
	msg, ok := ev.(model.EventStreamMsg)

	if !ok {
		fmt.Println("globalEventObserver Event fail cast to EventStreamMsg")
		return
	}

	if msg.Msg == constant.AddEventStream {
		es, ok := msg.Value.(*model.EventStream)
		if !ok {
			fmt.Printf("EventStreamMsg Value fail cast to *model.EventStream")
			return
		}

		for i, event := range es.Events {
			if i == 0 {
				ob.startEventNames = append(ob.startEventNames, event.Name)
			}
			ob.eventMap[event.Name] = event
		}
	} else if msg.Msg == constant.ResetEventStream {
		ob.startEventNames = []string{}
		ob.eventMap = map[string]*model.Event{}
	} else if msg.Msg == "start" {
		jsonMsg := model.JsonMsg{
			IsWrite: true,
			EventMsg: model.EventMsg{
				StartEventNames: ob.startEventNames,
				EventMap:        ob.eventMap,
			},
		}

		messagebus.SendMsg(constant.JsonFileObserverName, jsonMsg)

		//启动监控点击
		workMsg := model.WorkMsg{
			Msg: constant.StartState,
			EventMsg: model.EventMsg{
				StartEventNames: ob.startEventNames,
				EventMap:        ob.eventMap,
			},
		}
		messagebus.SendMsg(constant.WorkObserverName, workMsg)
	}
}
